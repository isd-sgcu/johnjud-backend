package auth

import (
	"context"

	authProto "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/auth/v1"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type AuthClientMock struct {
	mock.Mock
}

func (m *AuthClientMock) Validate(_ context.Context, in *authProto.ValidateRequest, _ ...grpc.CallOption) (*authProto.ValidateResponse, error) {
	args := m.Called(in)
	resp, _ := args.Get(0).(*authProto.ValidateResponse)
	err, _ := args.Get(1).(error)
	return resp, err
}

func (m *AuthClientMock) RefreshToken(_ context.Context, in *authProto.RefreshTokenRequest, _ ...grpc.CallOption) (*authProto.RefreshTokenResponse, error) {
	args := m.Called(in)
	resp, _ := args.Get(0).(*authProto.RefreshTokenResponse)
	err, _ := args.Get(1).(error)
	return resp, err
}

func (m *AuthClientMock) SignUp(_ context.Context, in *authProto.SignUpRequest, _ ...grpc.CallOption) (*authProto.SignUpResponse, error) {
	args := m.Called(in)
	resp, _ := args.Get(0).(*authProto.SignUpResponse)
	err, _ := args.Get(1).(error)
	return resp, err
}

func (m *AuthClientMock) SignIn(_ context.Context, in *authProto.SignInRequest, _ ...grpc.CallOption) (*authProto.SignInResponse, error) {
	args := m.Called(in)
	resp, _ := args.Get(0).(*authProto.SignInResponse)
	err, _ := args.Get(1).(error)
	return resp, err
}

func (m *AuthClientMock) SignOut(_ context.Context, in *authProto.SignOutRequest, _ ...grpc.CallOption) (*authProto.SignOutResponse, error) {
	args := m.Called(in)
	resp, _ := args.Get(0).(*authProto.SignOutResponse)
	err, _ := args.Get(1).(error)
	return resp, err
}

func (m *AuthClientMock) ForgotPassword(_ context.Context, in *authProto.ForgotPasswordRequest, _ ...grpc.CallOption) (*authProto.ForgotPasswordResponse, error) {
	args := m.Called(in)
	resp, _ := args.Get(0).(*authProto.ForgotPasswordResponse)
	err, _ := args.Get(1).(error)
	return resp, err
}

func (m *AuthClientMock) ResetPassword(_ context.Context, in *authProto.ResetPasswordRequest, _ ...grpc.CallOption) (*authProto.ResetPasswordResponse, error) {
	args := m.Called(in)
	resp, _ := args.Get(0).(*authProto.ResetPasswordResponse)
	err, _ := args.Get(1).(error)
	return resp, err
}
