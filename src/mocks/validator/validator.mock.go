// Code generated by MockGen. DO NOT EDIT.
// Source: ./src/app/validator/validator.go

// Package mock_validator is a generated GoMock package.
package mock_validator

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	dto "github.com/isd-sgcu/johnjud-gateway/src/app/dto"
)

// MockIDtoValidator is a mock of IDtoValidator interface.
type MockIDtoValidator struct {
	ctrl     *gomock.Controller
	recorder *MockIDtoValidatorMockRecorder
}

// MockIDtoValidatorMockRecorder is the mock recorder for MockIDtoValidator.
type MockIDtoValidatorMockRecorder struct {
	mock *MockIDtoValidator
}

// NewMockIDtoValidator creates a new mock instance.
func NewMockIDtoValidator(ctrl *gomock.Controller) *MockIDtoValidator {
	mock := &MockIDtoValidator{ctrl: ctrl}
	mock.recorder = &MockIDtoValidatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIDtoValidator) EXPECT() *MockIDtoValidatorMockRecorder {
	return m.recorder
}

// Validate mocks base method.
func (m *MockIDtoValidator) Validate(arg0 interface{}) []*dto.BadReqErrResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate", arg0)
	ret0, _ := ret[0].([]*dto.BadReqErrResponse)
	return ret0
}

// Validate indicates an expected call of Validate.
func (mr *MockIDtoValidatorMockRecorder) Validate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockIDtoValidator)(nil).Validate), arg0)
}
