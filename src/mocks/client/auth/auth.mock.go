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

func (m *AuthClientMock) Signup(_ context.Context, in *authProto.SignupRequest, _ ...grpc.CallOption) (*authProto.SignupResponse, error) {
	args := m.Called(in)
	resp, _ := args.Get(0).(*authProto.SignupResponse)
	err, _ := args.Get(1).(error)
	return resp, err
}

func (m *AuthClientMock) SignIn(_ context.Context, in *authProto.SignInRequest, _ ...grpc.CallOption) (*authProto.SignInResponse, error) {
	args := m.Called(in)
	resp, _ := args.Get(0).(*authProto.SignInResponse)
	err, _ := args.Get(1).(error)
	return resp, err
}
