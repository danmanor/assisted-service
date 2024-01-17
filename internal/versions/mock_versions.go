// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/openshift/assisted-service/internal/versions (interfaces: Handler)

// Package versions is a generated GoMock package.
package versions

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/openshift/assisted-service/models"
)

// MockHandler is a mock of Handler interface.
type MockHandler struct {
	ctrl     *gomock.Controller
	recorder *MockHandlerMockRecorder
}

// MockHandlerMockRecorder is the mock recorder for MockHandler.
type MockHandlerMockRecorder struct {
	mock *MockHandler
}

// NewMockHandler creates a new mock instance.
func NewMockHandler(ctrl *gomock.Controller) *MockHandler {
	mock := &MockHandler{ctrl: ctrl}
	mock.recorder = &MockHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHandler) EXPECT() *MockHandlerMockRecorder {
	return m.recorder
}

// GetDefaultReleaseImageByCPUArchitecture mocks base method.
func (m *MockHandler) GetDefaultReleaseImageByCPUArchitecture(arg0 string) (*models.ReleaseImage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDefaultReleaseImageByCPUArchitecture", arg0)
	ret0, _ := ret[0].(*models.ReleaseImage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDefaultReleaseImageByCPUArchitecture indicates an expected call of GetDefaultReleaseImageByCPUArchitecture.
func (mr *MockHandlerMockRecorder) GetDefaultReleaseImageByCPUArchitecture(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDefaultReleaseImageByCPUArchitecture", reflect.TypeOf((*MockHandler)(nil).GetDefaultReleaseImageByCPUArchitecture), arg0)
}

// GetMustGatherImages mocks base method.
func (m *MockHandler) GetMustGatherImages(arg0, arg1, arg2 string) (MustGatherVersion, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMustGatherImages", arg0, arg1, arg2)
	ret0, _ := ret[0].(MustGatherVersion)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMustGatherImages indicates an expected call of GetMustGatherImages.
func (mr *MockHandlerMockRecorder) GetMustGatherImages(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMustGatherImages", reflect.TypeOf((*MockHandler)(nil).GetMustGatherImages), arg0, arg1, arg2)
}

// GetReleaseImage mocks base method.
func (m *MockHandler) GetReleaseImage(arg0 context.Context, arg1, arg2, arg3 string) (*models.ReleaseImage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReleaseImage", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*models.ReleaseImage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReleaseImage indicates an expected call of GetReleaseImage.
func (mr *MockHandlerMockRecorder) GetReleaseImage(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReleaseImage", reflect.TypeOf((*MockHandler)(nil).GetReleaseImage), arg0, arg1, arg2, arg3)
}

// GetReleaseImageByURL mocks base method.
func (m *MockHandler) GetReleaseImageByURL(arg0 context.Context, arg1, arg2 string) (*models.ReleaseImage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReleaseImageByURL", arg0, arg1, arg2)
	ret0, _ := ret[0].(*models.ReleaseImage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReleaseImageByURL indicates an expected call of GetReleaseImageByURL.
func (mr *MockHandlerMockRecorder) GetReleaseImageByURL(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReleaseImageByURL", reflect.TypeOf((*MockHandler)(nil).GetReleaseImageByURL), arg0, arg1, arg2)
}

// ValidateReleaseImageForRHCOS mocks base method.
func (m *MockHandler) ValidateReleaseImageForRHCOS(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateReleaseImageForRHCOS", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateReleaseImageForRHCOS indicates an expected call of ValidateReleaseImageForRHCOS.
func (mr *MockHandlerMockRecorder) ValidateReleaseImageForRHCOS(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateReleaseImageForRHCOS", reflect.TypeOf((*MockHandler)(nil).ValidateReleaseImageForRHCOS), arg0, arg1)
}
