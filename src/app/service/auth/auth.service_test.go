package auth

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type AuthServiceTest struct {
	suite.Suite
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(AuthServiceTest))
}

func (t *AuthServiceTest) SetupTest() {

}

func (t *AuthServiceTest) TestSignupSuccess() {}

func (t *AuthServiceTest) TestSignupConflict() {}

func (t *AuthServiceTest) TestSignupInternalError() {}

func (t *AuthServiceTest) TestSignupUnavailableService() {}

func (t *AuthServiceTest) TestSignInSuccess() {}

func (t *AuthServiceTest) TestSignInForbidden() {}

func (t *AuthServiceTest) TestSignInInternalError() {}

func (t *AuthServiceTest) TestSignInUnavailableService() {}

func (t *AuthServiceTest) TestValidateSuccess() {}

func (t *AuthServiceTest) TestValidateUnauthorized() {}

func (t *AuthServiceTest) TestValidateInternalError() {}

func (t *AuthServiceTest) TestValidateUnavailableService() {}

func (t *AuthServiceTest) TestRefreshTokenUnauthorized() {}

func (t *AuthServiceTest) TestRefreshTokenInternalError() {}

func (t *AuthServiceTest) TestRefreshTokenUnavailableService() {}
