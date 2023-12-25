package auth

import (
	"github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	"github.com/isd-sgcu/johnjud-gateway/src/app/router"
	"github.com/isd-sgcu/johnjud-gateway/src/app/validator"
	"github.com/isd-sgcu/johnjud-gateway/src/pkg/service/auth"
	"github.com/isd-sgcu/johnjud-gateway/src/pkg/service/user"
	"net/http"
	"strings"
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
	request := &dto.SignupRequest{}
	err := c.Bind(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseErr{
			StatusCode: http.StatusBadRequest,
			Message:    constant.BindingRequestErrorMessage + err.Error(),
			Data:       nil,
		})
		return
	}

	if err := h.validate.Validate(request); err != nil {
		var errorMessage []string
		for _, reqErr := range err {
			errorMessage = append(errorMessage, reqErr.Message)
		}
		c.JSON(http.StatusBadRequest, dto.ResponseErr{
			StatusCode: http.StatusBadRequest,
			Message:    constant.InvalidRequestBodyMessage + strings.Join(errorMessage, ", "),
			Data:       nil,
		})
		return
	}

	response, respErr := h.service.Signup(request)
	if respErr != nil {
		c.JSON(respErr.StatusCode, respErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) Signin(c router.IContext) {

}

func (h *Handler) Signout(c router.IContext) {

}
