// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/user/user.service.go

// Package mock_user is a generated GoMock package.
package mock_user

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

// Delete mocks base method.
func (m *MockService) Delete(arg0 string) (*dto.DeleteUserResponse, *dto.ResponseErr) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(*dto.DeleteUserResponse)
	ret1, _ := ret[1].(*dto.ResponseErr)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockServiceMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockService)(nil).Delete), arg0)
}

// FindOne mocks base method.
func (m *MockService) FindOne(arg0 string) (*dto.User, *dto.ResponseErr) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOne", arg0)
	ret0, _ := ret[0].(*dto.User)
	ret1, _ := ret[1].(*dto.ResponseErr)
	return ret0, ret1
}

// FindOne indicates an expected call of FindOne.
func (mr *MockServiceMockRecorder) FindOne(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOne", reflect.TypeOf((*MockService)(nil).FindOne), arg0)
}

// Update mocks base method.
func (m *MockService) Update(arg0 string, arg1 *dto.UpdateUserRequest) (*dto.User, *dto.ResponseErr) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(*dto.User)
	ret1, _ := ret[1].(*dto.ResponseErr)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockServiceMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockService)(nil).Update), arg0, arg1)
}