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
	validate    validator.IDtoValidator
}

func NewHandler(service auth.Service, userService user.Service, validate validator.IDtoValidator) *Handler {
	return &Handler{service, userService, validate}
}

func (h *Handler) Validate(c router.IContext) {

}

func (h *Handler) RefreshToken(c router.IContext) {

}

func (h *Handler) Signup(c router.IContext) {
	// bind request
	// validate request
	// call authService.Signup
}

func (h *Handler) Signin(c router.IContext) {

}

func (h *Handler) Signout(c router.IContext) {

}
