package strategy

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type JwtStrategyTest struct {
	suite.Suite
	secret string
}

func TestJwtStrategy(t *testing.T) {
	suite.Run(t, new(JwtStrategyTest))
}

func (t *JwtStrategyTest) SetupTest() {
	secret := "testSecret"

	t.secret = secret
}

func (t *JwtStrategyTest) TestAuthDecodeSuccess() {
	token := &jwt.Token{
		Method: jwt.SigningMethodHS256,
	}
	expected := []byte(t.secret)

	jwtStrategy := NewJwtStrategy(t.secret)
	actual, err := jwtStrategy.AuthDecode(token)

	assert.Nil(t.T(), err)
	assert.Equal(t.T(), expected, actual)
}

func (t *JwtStrategyTest) TestAuthDecodeFailed() {
	token := &jwt.Token{
		Method: jwt.SigningMethodES256,
		Header: map[string]interface{}{
			"typ": "JWT",
			"alg": jwt.SigningMethodES256.Alg(),
		},
	}
	expected := errors.New(fmt.Sprintf("invalid token %v\n", token.Header["alg"]))

	jwtStrategy := NewJwtStrategy(t.secret)
	actual, err := jwtStrategy.AuthDecode(token)

	assert.Nil(t.T(), actual)
	assert.Equal(t.T(), expected.Error(), err.Error())
}
