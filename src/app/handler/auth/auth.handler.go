package auth

import (
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	"github.com/isd-sgcu/johnjud-gateway/src/app/router"
	"github.com/isd-sgcu/johnjud-gateway/src/app/validator"
	"github.com/isd-sgcu/johnjud-gateway/src/pkg/service/auth"
	user_proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/user/v1"
)

type Handler struct {
	service     auth.Service
	userService UserService
	validate    *validator.DtoValidator
}

type UserService interface {
	FindOne(string) (*user_proto.User, *dto.ResponseErr)
}

func NewHandler(service auth.Service, userService UserService, validate *validator.DtoValidator) *Handler {
	return &Handler{service, userService, validate}
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
