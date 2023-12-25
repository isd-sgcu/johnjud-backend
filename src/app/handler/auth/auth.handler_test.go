package auth

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type AuthHandlerTest struct {
	suite.Suite
}

func TestAuthHandler(t *testing.T) {
	suite.Run(t, new(AuthHandlerTest))
}

func (t *AuthHandlerTest) SetupTest() {

}

func (t *AuthHandlerTest) TestSignupSuccess() {}

func (t *AuthHandlerTest) TestSignupBindFailed() {}

func (t *AuthHandlerTest) TestSignupValidateFailed() {}

func (t *AuthHandlerTest) TestSignupServiceError() {}
