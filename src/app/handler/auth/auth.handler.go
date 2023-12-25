package auth

import (
	"github.com/isd-sgcu/johnjud-gateway/src/app/router"
	"github.com/isd-sgcu/johnjud-gateway/src/app/validator"
	"github.com/isd-sgcu/johnjud-gateway/src/pkg/service/auth"
	"github.com/isd-sgcu/johnjud-gateway/src/pkg/service/user"
)

type Handler struct {
	service     auth.Service
	userService user.Service
	validate    *validator.DtoValidator
}

func NewHandler(service auth.Service, userService user.Service, validate *validator.DtoValidator) *Handler {
	return &Handler{service, userService, validate}
}

func (h *Handler) Validate(c *router.FiberCtx) {

}

func (h *Handler) RefreshToken(c *router.FiberCtx) {

}

func (h *Handler) Signup(c *router.FiberCtx) {
	// bind request
	// validate request
	// call authService.Signup
}

func (h *Handler) Signin(c *router.FiberCtx) {

}

func (h *Handler) Signout(c *router.FiberCtx) {

}
