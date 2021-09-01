// Code generated by go-swagger; DO NOT EDIT.

package installer

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/openshift/assisted-service/models"
)

// V2GetClusterInstallConfigReader is a Reader for the V2GetClusterInstallConfig structure.
type V2GetClusterInstallConfigReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *V2GetClusterInstallConfigReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewV2GetClusterInstallConfigOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewV2GetClusterInstallConfigUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewV2GetClusterInstallConfigForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewV2GetClusterInstallConfigNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 405:
		result := NewV2GetClusterInstallConfigMethodNotAllowed()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewV2GetClusterInstallConfigInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewV2GetClusterInstallConfigOK creates a V2GetClusterInstallConfigOK with default headers values
func NewV2GetClusterInstallConfigOK() *V2GetClusterInstallConfigOK {
	return &V2GetClusterInstallConfigOK{}
}

/*V2GetClusterInstallConfigOK handles this case with default header values.

Success.
*/
type V2GetClusterInstallConfigOK struct {
	Payload string
}

func (o *V2GetClusterInstallConfigOK) Error() string {
	return fmt.Sprintf("[GET /v2/clusters/{cluster_id}/install-config][%d] v2GetClusterInstallConfigOK  %+v", 200, o.Payload)
}

func (o *V2GetClusterInstallConfigOK) GetPayload() string {
	return o.Payload
}

func (o *V2GetClusterInstallConfigOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewV2GetClusterInstallConfigUnauthorized creates a V2GetClusterInstallConfigUnauthorized with default headers values
func NewV2GetClusterInstallConfigUnauthorized() *V2GetClusterInstallConfigUnauthorized {
	return &V2GetClusterInstallConfigUnauthorized{}
}

/*V2GetClusterInstallConfigUnauthorized handles this case with default header values.

Unauthorized.
*/
type V2GetClusterInstallConfigUnauthorized struct {
	Payload *models.InfraError
}

func (o *V2GetClusterInstallConfigUnauthorized) Error() string {
	return fmt.Sprintf("[GET /v2/clusters/{cluster_id}/install-config][%d] v2GetClusterInstallConfigUnauthorized  %+v", 401, o.Payload)
}

func (o *V2GetClusterInstallConfigUnauthorized) GetPayload() *models.InfraError {
	return o.Payload
}

func (o *V2GetClusterInstallConfigUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.InfraError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewV2GetClusterInstallConfigForbidden creates a V2GetClusterInstallConfigForbidden with default headers values
func NewV2GetClusterInstallConfigForbidden() *V2GetClusterInstallConfigForbidden {
	return &V2GetClusterInstallConfigForbidden{}
}

/*V2GetClusterInstallConfigForbidden handles this case with default header values.

Forbidden.
*/
type V2GetClusterInstallConfigForbidden struct {
	Payload *models.InfraError
}

func (o *V2GetClusterInstallConfigForbidden) Error() string {
	return fmt.Sprintf("[GET /v2/clusters/{cluster_id}/install-config][%d] v2GetClusterInstallConfigForbidden  %+v", 403, o.Payload)
}

func (o *V2GetClusterInstallConfigForbidden) GetPayload() *models.InfraError {
	return o.Payload
}

func (o *V2GetClusterInstallConfigForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.InfraError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewV2GetClusterInstallConfigNotFound creates a V2GetClusterInstallConfigNotFound with default headers values
func NewV2GetClusterInstallConfigNotFound() *V2GetClusterInstallConfigNotFound {
	return &V2GetClusterInstallConfigNotFound{}
}

/*V2GetClusterInstallConfigNotFound handles this case with default header values.

Error.
*/
type V2GetClusterInstallConfigNotFound struct {
	Payload *models.Error
}

func (o *V2GetClusterInstallConfigNotFound) Error() string {
	return fmt.Sprintf("[GET /v2/clusters/{cluster_id}/install-config][%d] v2GetClusterInstallConfigNotFound  %+v", 404, o.Payload)
}

func (o *V2GetClusterInstallConfigNotFound) GetPayload() *models.Error {
	return o.Payload
}

func (o *V2GetClusterInstallConfigNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewV2GetClusterInstallConfigMethodNotAllowed creates a V2GetClusterInstallConfigMethodNotAllowed with default headers values
func NewV2GetClusterInstallConfigMethodNotAllowed() *V2GetClusterInstallConfigMethodNotAllowed {
	return &V2GetClusterInstallConfigMethodNotAllowed{}
}

/*V2GetClusterInstallConfigMethodNotAllowed handles this case with default header values.

Method Not Allowed.
*/
type V2GetClusterInstallConfigMethodNotAllowed struct {
	Payload *models.Error
}

func (o *V2GetClusterInstallConfigMethodNotAllowed) Error() string {
	return fmt.Sprintf("[GET /v2/clusters/{cluster_id}/install-config][%d] v2GetClusterInstallConfigMethodNotAllowed  %+v", 405, o.Payload)
}

func (o *V2GetClusterInstallConfigMethodNotAllowed) GetPayload() *models.Error {
	return o.Payload
}

func (o *V2GetClusterInstallConfigMethodNotAllowed) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewV2GetClusterInstallConfigInternalServerError creates a V2GetClusterInstallConfigInternalServerError with default headers values
func NewV2GetClusterInstallConfigInternalServerError() *V2GetClusterInstallConfigInternalServerError {
	return &V2GetClusterInstallConfigInternalServerError{}
}

/*V2GetClusterInstallConfigInternalServerError handles this case with default header values.

Error.
*/
type V2GetClusterInstallConfigInternalServerError struct {
	Payload *models.Error
}

func (o *V2GetClusterInstallConfigInternalServerError) Error() string {
	return fmt.Sprintf("[GET /v2/clusters/{cluster_id}/install-config][%d] v2GetClusterInstallConfigInternalServerError  %+v", 500, o.Payload)
}

func (o *V2GetClusterInstallConfigInternalServerError) GetPayload() *models.Error {
	return o.Payload
}

func (o *V2GetClusterInstallConfigInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
