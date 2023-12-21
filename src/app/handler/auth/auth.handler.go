package auth

import (
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	"github.com/isd-sgcu/johnjud-gateway/src/app/router"
	"github.com/isd-sgcu/johnjud-gateway/src/app/validator"
	auth_proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/auth/v1"
	user_proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/user/v1"
)

type Handler struct {
	service    Service
	usrService UserService
	validate   *validator.DtoValidator
}

type Service interface {
	Signup(*dto.Signup) (*auth_proto.Credential, *dto.ResponseErr)
	Signin(*dto.Signin) (*auth_proto.Credential, *dto.ResponseErr)
	Signout(string) (bool, *dto.ResponseErr)
	Validate(string) (*dto.TokenPayloadAuth, *dto.ResponseErr)
	RefreshToken(string) (*auth_proto.Credential, *dto.ResponseErr)
}

type UserService interface {
	FindOne(string) (*user_proto.User, *dto.ResponseErr)
}

func NewHandler(service Service, usrService UserService, validate *validator.DtoValidator) *Handler {
	return &Handler{service, usrService, validate}
}

func (h *Handler) Validate(c *router.FiberCtx) {

}

func (h *Handler) RefreshToken(c *router.FiberCtx) {

}

func (h *Handler) Signup(c *router.FiberCtx) {

}

func (h *Handler) Signin(c *router.FiberCtx) {

}

func (h *Handler) Signout(c *router.FiberCtx) {

}
