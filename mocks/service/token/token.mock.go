// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/auth/token/token.service.go

// Package mock_token is a generated GoMock package.
package mock_token

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	constant "github.com/isd-sgcu/johnjud-backend/constant"
	dto "github.com/isd-sgcu/johnjud-backend/internal/dto"
	v1 "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/auth/v1"
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

// CreateCredential mocks base method.
func (m *MockService) CreateCredential(userId string, role constant.Role, authSessionId string) (*v1.Credential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCredential", userId, role, authSessionId)
	ret0, _ := ret[0].(*v1.Credential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCredential indicates an expected call of CreateCredential.
func (mr *MockServiceMockRecorder) CreateCredential(userId, role, authSessionId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCredential", reflect.TypeOf((*MockService)(nil).CreateCredential), userId, role, authSessionId)
}

// CreateRefreshToken mocks base method.
func (m *MockService) CreateRefreshToken() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRefreshToken")
	ret0, _ := ret[0].(string)
	return ret0
}

// CreateRefreshToken indicates an expected call of CreateRefreshToken.
func (mr *MockServiceMockRecorder) CreateRefreshToken() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRefreshToken", reflect.TypeOf((*MockService)(nil).CreateRefreshToken))
}

// CreateResetPasswordToken mocks base method.
func (m *MockService) CreateResetPasswordToken(userId string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateResetPasswordToken", userId)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateResetPasswordToken indicates an expected call of CreateResetPasswordToken.
func (mr *MockServiceMockRecorder) CreateResetPasswordToken(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateResetPasswordToken", reflect.TypeOf((*MockService)(nil).CreateResetPasswordToken), userId)
}

// FindRefreshTokenCache mocks base method.
func (m *MockService) FindRefreshTokenCache(refreshToken string) (*dto.RefreshTokenCache, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindRefreshTokenCache", refreshToken)
	ret0, _ := ret[0].(*dto.RefreshTokenCache)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindRefreshTokenCache indicates an expected call of FindRefreshTokenCache.
func (mr *MockServiceMockRecorder) FindRefreshTokenCache(refreshToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindRefreshTokenCache", reflect.TypeOf((*MockService)(nil).FindRefreshTokenCache), refreshToken)
}

// FindResetPasswordToken mocks base method.
func (m *MockService) FindResetPasswordToken(token string) (*dto.ResetPasswordTokenCache, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindResetPasswordToken", token)
	ret0, _ := ret[0].(*dto.ResetPasswordTokenCache)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindResetPasswordToken indicates an expected call of FindResetPasswordToken.
func (mr *MockServiceMockRecorder) FindResetPasswordToken(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindResetPasswordToken", reflect.TypeOf((*MockService)(nil).FindResetPasswordToken), token)
}

// RemoveAccessTokenCache mocks base method.
func (m *MockService) RemoveAccessTokenCache(authSessionId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveAccessTokenCache", authSessionId)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveAccessTokenCache indicates an expected call of RemoveAccessTokenCache.
func (mr *MockServiceMockRecorder) RemoveAccessTokenCache(authSessionId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveAccessTokenCache", reflect.TypeOf((*MockService)(nil).RemoveAccessTokenCache), authSessionId)
}

// RemoveRefreshTokenCache mocks base method.
func (m *MockService) RemoveRefreshTokenCache(refreshToken string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveRefreshTokenCache", refreshToken)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveRefreshTokenCache indicates an expected call of RemoveRefreshTokenCache.
func (mr *MockServiceMockRecorder) RemoveRefreshTokenCache(refreshToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveRefreshTokenCache", reflect.TypeOf((*MockService)(nil).RemoveRefreshTokenCache), refreshToken)
}

// RemoveResetPasswordToken mocks base method.
func (m *MockService) RemoveResetPasswordToken(token string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveResetPasswordToken", token)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveResetPasswordToken indicates an expected call of RemoveResetPasswordToken.
func (mr *MockServiceMockRecorder) RemoveResetPasswordToken(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveResetPasswordToken", reflect.TypeOf((*MockService)(nil).RemoveResetPasswordToken), token)
}

// Validate mocks base method.
func (m *MockService) Validate(token string) (*dto.UserCredential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate", token)
	ret0, _ := ret[0].(*dto.UserCredential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Validate indicates an expected call of Validate.
func (mr *MockServiceMockRecorder) Validate(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockService)(nil).Validate), token)
}
