package staticnetworkconfig

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/nmstate/nmstate/rust/src/go/nmstate"
	"github.com/openshift/assisted-service/models"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/util/json"
	yamlconvertor "sigs.k8s.io/yaml"
)

type StaticNetworkConfigData struct {
	FilePath     string
	FileContents string
}

type networkInterface struct {
	Name          string `yaml:"name"`
	InterfaceType string `yaml:"type"`
	Vlan          *vlan  `yaml:"vlan,omitempty"`
	Identifier    string `yaml:"identifier,omitempty"`
}

type vlan struct {
	BaseIface string `yaml:"base-iface"`
}
type interfaceConfig struct {
	Interfaces []networkInterface `yaml:"interfaces"`
}

//go:generate mockgen -source=generator.go -package=staticnetworkconfig -destination=mock_generator.go
type StaticNetworkConfig interface {
	GenerateStaticNetworkConfigData(ctx context.Context, hostsYAMLS string) ([]StaticNetworkConfigData, error)
	GenerateStaticNetworkConfigDataYAML(staticNetworkConfigStr string) ([]StaticNetworkConfigData, error)
	FormatStaticNetworkConfigForDB(staticNetworkConfig []*models.HostStaticNetworkConfig) (string, error)
	ValidateStaticConfigParamsYAML(staticNetworkConfig []*models.HostStaticNetworkConfig, ocpVersion, arch, installerInvoker string) error
}

type StaticNetworkConfigGenerator struct {
	log     logrus.FieldLogger
	nmstate *nmstate.Nmstate
}

func New(log logrus.FieldLogger) StaticNetworkConfig {
	return &StaticNetworkConfigGenerator{
		log:     log,
		nmstate: nmstate.New(),
	}
}

func (s *StaticNetworkConfigGenerator) injectNMPolicyCaptures(hostConfig *models.HostStaticNetworkConfig) (string, error) {
	interfacesWithCaptures := hostConfig.NetworkYaml
	var config interfaceConfig

	err := yaml.Unmarshal([]byte(interfacesWithCaptures), &config)
	if err != nil {
		s.log.WithError(err).Errorf("error unmarshalling JSON: %v", err)
		return "", err
	}

	nameToType := make(map[string]string)
	for _, iface := range config.Interfaces {
		if iface.InterfaceType == "vlan" {
			nameToType[iface.Vlan.BaseIface] = iface.InterfaceType
		} else {
			nameToType[iface.Name] = iface.InterfaceType
		}
	}

	var captureSection []string

	for j, mac := range hostConfig.MacInterfaceMap {
		macAddress := mac.MacAddress
		interfaceName := mac.LogicalNicName

		// Generate capture section
		captureSection = append(captureSection, fmt.Sprintf("  iface%d: interfaces.mac-address == \"%s\"", j, strings.ToUpper(macAddress)))

		// Replace logical names with capture values.
		var modifiedInterfaceConfig string
		// VLAN is a special case where the base-iface refers to the logical name rather than the VLAN's name.
		if val, exists := nameToType[mac.LogicalNicName]; exists && val == "vlan" {
			re := regexp.MustCompile(fmt.Sprintf("base-iface: %s", interfaceName))
			modifiedInterfaceConfig = re.ReplaceAllString(interfacesWithCaptures, fmt.Sprintf(`base-iface: "{{ capture.iface%d.interfaces.0.name }}"`, j))
		} else {
			re := regexp.MustCompile(interfaceName)
			modifiedInterfaceConfig = re.ReplaceAllString(interfacesWithCaptures, fmt.Sprintf(`"{{ capture.iface%d.interfaces.0.name }}"`, j))

		}

		// Add indentation
		var indentedConfig []string
		for _, line := range strings.Split(modifiedInterfaceConfig, "\n") {
			// Check if line is non-empty and add indentation
			if strings.TrimSpace(line) != "" {
				indentedConfig = append(indentedConfig, "  "+line)
			}
		}
		interfacesWithCaptures = strings.Join(indentedConfig, "\n")
	}
	// Generate desiredState section
	desiredStateSection := "\ndesiredState:\n"

	// Combine the sections
	finalOutput := "capture:\n" + strings.Join(captureSection, "\n") + desiredStateSection + interfacesWithCaptures
	return finalOutput, nil
}

func (s *StaticNetworkConfigGenerator) GenerateStaticNetworkConfigDataYAML(staticNetworkConfigStr string) ([]StaticNetworkConfigData, error) {
	var filesList []StaticNetworkConfigData
	staticNetworkConfig, err := s.decodeStaticNetworkConfig(staticNetworkConfigStr)
	if err != nil {
		s.log.WithError(err).Errorf("Failed to decode static network config")
		return nil, err
	}

	for i, hostConfig := range staticNetworkConfig {
		// We currently use this to validate the YAML file only
		_, err = s.generateConfiguration(hostConfig.NetworkYaml)
		if err != nil {
			return nil, err
		}

		nmpolicy, err := s.injectNMPolicyCaptures(hostConfig)
		if err != nil {
			return nil, err
		}

		macInterfaceMapping := s.formatMacInterfaceMap(hostConfig.MacInterfaceMap)
		mapConfigData := StaticNetworkConfigData{
			FilePath:     filepath.Join(fmt.Sprintf("host%d", i), "mac_interface.ini"),
			FileContents: macInterfaceMapping,
		}
		filesList = append(filesList, mapConfigData)

		newFile := StaticNetworkConfigData{
			FilePath:     filepath.Join(filepath.Join(fmt.Sprintf("host%d", i), fmt.Sprintf("ymlFile%d.yml", i))),
			FileContents: nmpolicy,
		}
		filesList = append(filesList, newFile)
	}
	return filesList, nil
}

func (s *StaticNetworkConfigGenerator) GenerateStaticNetworkConfigData(ctx context.Context, staticNetworkConfigStr string) ([]StaticNetworkConfigData, error) {

	staticNetworkConfig, err := s.decodeStaticNetworkConfig(staticNetworkConfigStr)
	if err != nil {
		s.log.WithError(err).Errorf("Failed to decode static network config")
		return nil, err
	}
	s.log.Infof("Start configuring static network for %d hosts", len(staticNetworkConfig))
	filesList := []StaticNetworkConfigData{}
	for i, hostConfig := range staticNetworkConfig {
		hostFileList, err := s.generateHostStaticNetworkConfigData(hostConfig, fmt.Sprintf("host%d", i))
		if err != nil {
			err = errors.Wrapf(err, "failed to create static config for host %d", i)
			s.log.Error(err)
			return nil, err
		}
		filesList = append(filesList, hostFileList...)
	}
	return filesList, nil
}

func (s *StaticNetworkConfigGenerator) generateHostStaticNetworkConfigData(hostConfig *models.HostStaticNetworkConfig, hostDir string) ([]StaticNetworkConfigData, error) {
	hostYAML := hostConfig.NetworkYaml
	macInterfaceMapping := s.formatMacInterfaceMap(hostConfig.MacInterfaceMap)
	result, err := s.generateConfiguration(hostYAML)
	if err != nil {
		return nil, err
	}
	filesList, err := s.createNMConnectionFiles(result, hostDir)
	if err != nil {
		s.log.WithError(err).Errorf("failed to create NM connection files")
		return nil, err
	}
	mapConfigData := StaticNetworkConfigData{
		FilePath:     filepath.Join(hostDir, "mac_interface.ini"),
		FileContents: macInterfaceMapping,
	}
	filesList = append(filesList, mapConfigData)
	return filesList, nil
}

func (s *StaticNetworkConfigGenerator) generateConfiguration(hostYAML string) (string, error) {
	if hostYAML == "" {
		return "", errors.New("cannot generate configuration with an empty host YAML")
	}
	hostJSON, err := yamlconvertor.YAMLToJSON([]byte(hostYAML))
	if err != nil {
		return "", err
	}
	stdout, err := s.nmstate.GenerateConfiguration(string(hostJSON))
	if err != nil {
		s.log.WithError(err).Errorf("nmstate GenerateConfiguration failed, input yaml <%s>", hostYAML)
		return "", fmt.Errorf("nmstate GenerateConfiguration failed, error: %s", err.Error())
	}

	return stdout, nil
}

// create NMConnectionFiles formats the nmstate output into a list of file data
// Nothing is written to the local filesystem
func (s *StaticNetworkConfigGenerator) createNMConnectionFiles(nmstateOutput, hostDir string) ([]StaticNetworkConfigData, error) {
	var hostNMConnections map[string]interface{}
	err := yaml.Unmarshal([]byte(nmstateOutput), &hostNMConnections)
	if err != nil {
		s.log.WithError(err).Errorf("Failed to unmarshal nmstate output")
		return nil, err
	}
	if _, found := hostNMConnections["NetworkManager"]; !found {
		return nil, errors.Errorf("nmstate generated an empty NetworkManager config file content")
	}
	filesList := []StaticNetworkConfigData{}
	connectionsList := hostNMConnections["NetworkManager"].([]interface{})
	if len(connectionsList) == 0 {
		return nil, errors.Errorf("nmstate generated an empty NetworkManager config file content")
	}
	for _, connection := range connectionsList {
		connectionElems := connection.([]interface{})
		fileName := connectionElems[0].(string)
		fileContents, err := s.formatNMConnection(connectionElems[1].(string))
		if err != nil {
			return nil, err
		}
		s.log.Infof("Adding NMConnection file <%s>", fileName)
		newFile := StaticNetworkConfigData{
			FilePath:     filepath.Join(hostDir, fileName),
			FileContents: fileContents,
		}
		filesList = append(filesList, newFile)
	}
	return filesList, nil
}

func (s *StaticNetworkConfigGenerator) formatNMConnection(nmConnection string) (string, error) {
	ini.PrettyFormat = false
	cfg, err := ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, []byte(nmConnection))
	if err != nil {
		s.log.WithError(err).Errorf("Failed to load the ini format string %s", nmConnection)
		return "", err
	}
	connectionSection := cfg.Section("connection")
	_, err = connectionSection.NewKey("autoconnect", "true")
	if err != nil {
		s.log.WithError(err).Errorf("Failed to add autoconnect key to section connection")
		return "", err
	}
	_, err = connectionSection.NewKey("autoconnect-priority", "1")
	if err != nil {
		s.log.WithError(err).Errorf("Failed to add autoconnect-priority key to section connection")
		return "", err
	}

	buf := new(bytes.Buffer)
	_, err = cfg.WriteTo(buf)
	if err != nil {
		s.log.WithError(err).Errorf("Failed to output nmconnection ini file to buffer")
		return "", err
	}
	return buf.String(), nil
}

func (s *StaticNetworkConfigGenerator) validateInterfaceNamesExistenceYAML(macInterfaceMap models.MacInterfaceMap, networksYaml, ocpVersion, arch, installerInvoker string) error {
	interfaceNames := lo.Map(macInterfaceMap, func(m *models.MacInterfaceMapItems0, _ int) string { return m.LogicalNicName })

	var config map[string]interface{}

	// Unmarshal the JSON string into the config struct
	err := yaml.Unmarshal([]byte(networksYaml), &config)
	if err != nil {
		s.log.WithError(err).Errorf("Error unmarshalling JSON")
		return err
	}

	// See 'man systemd.net-naming-scheme' for interface naming protocol
	predicatableIfaceNamePattern := regexp.MustCompile("^en[PsvxXbucaipod]")
	Interfaces := config["interfaces"].([]interface{})
	interfaceWithmacIdentifier := make(map[string]struct{})

	for _, iface := range Interfaces {
		nic := iface.(map[interface{}]interface{})
		interfaceName, exists := nic["name"]
		if !exists {
			return errors.Errorf("interface name not found in networks configuration")
		}

		identifier, exists := nic["identifier"]
		if exists && identifier == "mac-address" {
			interfaceWithmacIdentifier[interfaceName.(string)] = struct{}{}
		}
	}

	for _, iface := range Interfaces {
		nic := iface.(map[interface{}]interface{})
		interfaceName, exists := nic["name"]
		if !exists {
			return errors.Errorf("interface name not found in networks configuration")
		}

		identifier, exists := nic["identifier"]
		isMacAddressIdentifier := exists && identifier == "mac-address"
		isVersionOK, err := NMStatectlServiceSupported(ocpVersion, arch)
		if err != nil {
			return err
		}

		// TODO: This is a temporary validation. It should be removed once the mac-identifier bug in nmstate is resolved in version 2.2.35.
		if installerInvoker != "agent-installer" && isVersionOK && isMacAddressIdentifier {
			return errors.Errorf("Mac-address identifier is not supported")
		}

		interfaceType := nic["type"]
		switch interfaceType {
		case "802-3-ethernet", "ethernet":
			if !lo.Contains(interfaceNames, interfaceName.(string)) {
				if isMacAddressIdentifier {
					s.log.Infof("Interface %s has no mac-interface mapping but has a mac-identifier", interfaceName)
				} else if !predicatableIfaceNamePattern.MatchString(interfaceName.(string)) {
					return errors.Errorf("mac-interface mapping for interface %s is missing and not a physical interface", interfaceName)
				} else {
					s.log.Infof("Interface %s has no mac-interface mapping but matches a physical interface", interfaceName)
				}
			}
		case "bond":
			if val, exists := nic["link-aggregation"].(map[interface{}]interface{}); exists {
				if ports, exists := val["port"].([]interface{}); exists {
					for _, port := range ports {
						if !lo.Contains(interfaceNames, port.(string)) {
							if _, exists := interfaceWithmacIdentifier[port.(string)]; exists {
								continue
							} else if !predicatableIfaceNamePattern.MatchString(port.(string)) {
								return errors.Errorf("mac-interface mapping for interface %s is missing and not a physical interface", interfaceName)
							} else {
								s.log.Infof("Interface %s has no mac-interface mapping but matches a physical interface", interfaceName)
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func (s *StaticNetworkConfigGenerator) ValidateStaticConfigParamsYAML(staticNetworkConfig []*models.HostStaticNetworkConfig, ocpVersion, arch, installInvoker string) error {
	var err *multierror.Error
	for i, hostConfig := range staticNetworkConfig {
		err = multierror.Append(err, s.validateMacInterfaceName(i, hostConfig.MacInterfaceMap))
		_, validateErr := s.generateConfiguration(hostConfig.NetworkYaml)
		if validateErr != nil {
			err = multierror.Append(err, fmt.Errorf("failed to validate network yaml for host %d, %s", i, validateErr))
		}
		err = multierror.Append(err, s.validateInterfaceNamesExistenceYAML(hostConfig.MacInterfaceMap, hostConfig.NetworkYaml, ocpVersion, arch, installInvoker))
	}
	return err.ErrorOrNil()
}

func (s *StaticNetworkConfigGenerator) validateMacInterfaceName(hostIdx int, macInterfaceMap models.MacInterfaceMap) error {
	if len(macInterfaceMap) == 0 {
		return fmt.Errorf("at least one interface for host %d must be provided", hostIdx)
	}
	interfaceCheck := make(map[string]struct{}, len(macInterfaceMap))
	macCheck := make(map[string]struct{}, len(macInterfaceMap))
	for _, macInterface := range macInterfaceMap {
		interfaceCheck[macInterface.LogicalNicName] = struct{}{}
		macCheck[macInterface.MacAddress] = struct{}{}
	}
	if len(interfaceCheck) < len(macInterfaceMap) || len(macCheck) < len(macInterfaceMap) {
		return fmt.Errorf("MACs and Interfaces for host %d must be unique", hostIdx)
	}
	return nil
}

func compareMapInterfaces(intf1, intf2 *models.MacInterfaceMapItems0) bool {
	if intf1.LogicalNicName != intf2.LogicalNicName {
		return intf1.LogicalNicName < intf2.LogicalNicName
	}
	return intf1.MacAddress < intf2.MacAddress
}

func compareMacInterfaceMaps(map1, map2 models.MacInterfaceMap) bool {
	if len(map1) != len(map2) {
		return len(map1) < len(map2)
	}
	for i := range map1 {
		less := compareMapInterfaces(map1[i], map2[i])
		greater := compareMapInterfaces(map2[i], map1[i])
		if less || greater {
			return less
		}
	}
	return false
}

func sortStaticNetworkConfig(staticNetworkConfig []*models.HostStaticNetworkConfig) {
	for i := range staticNetworkConfig {
		item := staticNetworkConfig[i]
		sort.SliceStable(item.MacInterfaceMap, func(i, j int) bool {
			return compareMapInterfaces(item.MacInterfaceMap[i], item.MacInterfaceMap[j])
		})
	}
	sort.SliceStable(staticNetworkConfig, func(i, j int) bool {
		hostConfig1 := staticNetworkConfig[i]
		hostConfig2 := staticNetworkConfig[j]
		if hostConfig1.NetworkYaml != hostConfig2.NetworkYaml {
			return hostConfig1.NetworkYaml < hostConfig2.NetworkYaml
		}
		return compareMacInterfaceMaps(hostConfig1.MacInterfaceMap, hostConfig2.MacInterfaceMap)
	})
}

func (s *StaticNetworkConfigGenerator) FormatStaticNetworkConfigForDB(staticNetworkConfig []*models.HostStaticNetworkConfig) (string, error) {
	if len(staticNetworkConfig) == 0 {
		return "", nil
	}
	sortStaticNetworkConfig(staticNetworkConfig)
	b, err := json.Marshal(&staticNetworkConfig)
	if err != nil {
		return "", errors.Wrap(err, "Failed to JSON Marshal static network config")
	}
	return string(b), nil
}

func (s *StaticNetworkConfigGenerator) decodeStaticNetworkConfig(staticNetworkConfigStr string) (staticNetworkConfig []*models.HostStaticNetworkConfig, err error) {
	if staticNetworkConfigStr == "" {
		return
	}
	err = json.Unmarshal([]byte(staticNetworkConfigStr), &staticNetworkConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to JSON Unmarshal static network config %s", staticNetworkConfigStr)
	}
	return
}

func (s *StaticNetworkConfigGenerator) formatMacInterfaceMap(macInterfaceMap models.MacInterfaceMap) string {
	lines := make([]string, len(macInterfaceMap))
	for i, entry := range macInterfaceMap {
		lines[i] = fmt.Sprintf("%s=%s", entry.MacAddress, entry.LogicalNicName)
	}
	sort.Strings(lines)
	return strings.Join(lines, "\n")
}

func GenerateStaticNetworkConfigArchive(files []StaticNetworkConfigData) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	w := tar.NewWriter(buffer)
	for _, file := range files {
		path := filepath.Join("/etc/assisted/network", file.FilePath)
		content := file.FileContents

		// add the file content
		hdr := &tar.Header{
			Name: path,
			Mode: 0600,
			Size: int64(len(content)),
		}
		if err := w.WriteHeader(hdr); err != nil {
			return nil, err
		}
		if _, err := w.Write([]byte(content)); err != nil {
			return nil, err
		}
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buffer, nil
}
