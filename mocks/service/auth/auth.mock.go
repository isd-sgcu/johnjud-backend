// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/auth/auth.service.go

// Package mock_auth is a generated GoMock package.
package mock_auth

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	dto "github.com/isd-sgcu/johnjud-gateway/internal/dto"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// ForgotPassword mocks base method.
func (m *MockService) ForgotPassword(request *dto.ForgotPasswordRequest) (*dto.ForgotPasswordResponse, *dto.ResponseErr) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForgotPassword", request)
	ret0, _ := ret[0].(*dto.ForgotPasswordResponse)
	ret1, _ := ret[1].(*dto.ResponseErr)
	return ret0, ret1
}

// ForgotPassword indicates an expected call of ForgotPassword.
func (mr *MockServiceMockRecorder) ForgotPassword(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForgotPassword", reflect.TypeOf((*MockService)(nil).ForgotPassword), request)
}

// RefreshToken mocks base method.
func (m *MockService) RefreshToken(request *dto.RefreshTokenRequest) (*dto.Credential, *dto.ResponseErr) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshToken", request)
	ret0, _ := ret[0].(*dto.Credential)
	ret1, _ := ret[1].(*dto.ResponseErr)
	return ret0, ret1
}

// RefreshToken indicates an expected call of RefreshToken.
func (mr *MockServiceMockRecorder) RefreshToken(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshToken", reflect.TypeOf((*MockService)(nil).RefreshToken), request)
}

// ResetPassword mocks base method.
func (m *MockService) ResetPassword(request *dto.ResetPasswordRequest) (*dto.ResetPasswordResponse, *dto.ResponseErr) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetPassword", request)
	ret0, _ := ret[0].(*dto.ResetPasswordResponse)
	ret1, _ := ret[1].(*dto.ResponseErr)
	return ret0, ret1
}

// ResetPassword indicates an expected call of ResetPassword.
func (mr *MockServiceMockRecorder) ResetPassword(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetPassword", reflect.TypeOf((*MockService)(nil).ResetPassword), request)
}

// SignIn mocks base method.
func (m *MockService) SignIn(request *dto.SignInRequest) (*dto.Credential, *dto.ResponseErr) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignIn", request)
	ret0, _ := ret[0].(*dto.Credential)
	ret1, _ := ret[1].(*dto.ResponseErr)
	return ret0, ret1
}

// SignIn indicates an expected call of SignIn.
func (mr *MockServiceMockRecorder) SignIn(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignIn", reflect.TypeOf((*MockService)(nil).SignIn), request)
}

// SignOut mocks base method.
func (m *MockService) SignOut(accessToken string) (*dto.SignOutResponse, *dto.ResponseErr) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignOut", accessToken)
	ret0, _ := ret[0].(*dto.SignOutResponse)
	ret1, _ := ret[1].(*dto.ResponseErr)
	return ret0, ret1
}

// SignOut indicates an expected call of SignOut.
func (mr *MockServiceMockRecorder) SignOut(accessToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignOut", reflect.TypeOf((*MockService)(nil).SignOut), accessToken)
}

// Signup mocks base method.
func (m *MockService) Signup(request *dto.SignupRequest) (*dto.SignupResponse, *dto.ResponseErr) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Signup", request)
	ret0, _ := ret[0].(*dto.SignupResponse)
	ret1, _ := ret[1].(*dto.ResponseErr)
	return ret0, ret1
}

// Signup indicates an expected call of Signup.
func (mr *MockServiceMockRecorder) Signup(request interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Signup", reflect.TypeOf((*MockService)(nil).Signup), request)
}

// Validate mocks base method.
func (m *MockService) Validate(refreshToken string) (*dto.TokenPayloadAuth, *dto.ResponseErr) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate", refreshToken)
	ret0, _ := ret[0].(*dto.TokenPayloadAuth)
	ret1, _ := ret[1].(*dto.ResponseErr)
	return ret0, ret1
}

// Validate indicates an expected call of Validate.
func (mr *MockServiceMockRecorder) Validate(refreshToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockService)(nil).Validate), refreshToken)
}
