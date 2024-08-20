package auth

import (
	"net/http"
	"strings"

	"github.com/isd-sgcu/johnjud-gateway/constant"
	"github.com/isd-sgcu/johnjud-gateway/internal/dto"
	"github.com/isd-sgcu/johnjud-gateway/internal/pkg/service/user"
	"github.com/isd-sgcu/johnjud-gateway/internal/router"
	"github.com/isd-sgcu/johnjud-gateway/internal/validator"
)

type handlerImpl struct {
	service     Service
	userService user.Service
	validate    validator.IDtoValidator
}

func NewHandler(service Service, userService user.Service, validate validator.IDtoValidator) *handlerImpl {
	return &handlerImpl{service, userService, validate}
}

// Signup is a function that create user in database
// @Summary Signup user
// @Description Return the data of user if successfully
// @Param signup body dto.SignupRequest true "signup request dto"
// @Tags auth
// @Accept json
// @Produce json
// @Success 201 {object} dto.SignupResponse
// @Failure 400 {object} dto.ResponseBadRequestErr "Invalid request body"
// @Failure 409 {object} dto.ResponseConflictErr "Duplicate email"
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Router /v1/auth/signup [post]
func (h *handlerImpl) Signup(c router.IContext) {
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

// SignIn is a function that authenticate user with email and password
// @Summary Sign in user
// @Description Return the credential of user including access token and refresh token
// @Param signIn body dto.SignInRequest true "signIn request dto"
// @Tags auth
// @Accept json
// @Produce json
// @Success 201 {object} dto.Credential
// @Failure 400 {object} dto.ResponseBadRequestErr "Invalid request body"
// @Failure 403 {object} dto.ResponseForbiddenErr "Incorrect email or password"
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Router /v1/auth/signin [post]
func (h *handlerImpl) SignIn(c router.IContext) {
	request := &dto.SignInRequest{}
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

	response, respErr := h.service.SignIn(request)
	if respErr != nil {
		c.JSON(respErr.StatusCode, respErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

// SignOut is a function that remove token and auth session of user
// @Summary Sign out user
// @Description Return the bool value of success
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.SignOutResponse
// @Failure 401 {object} dto.ResponseUnauthorizedErr "Invalid token"
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Router /v1/auth/signout [post]
func (h *handlerImpl) SignOut(c router.IContext) {
	token := c.Token()

	response, respErr := h.service.SignOut(token)
	if respErr != nil {
		c.JSON(respErr.StatusCode, respErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken is a function to redeem new access token and refresh token
// @Summary Refresh token
// @Description Return the credential
// @Param request body dto.RefreshTokenRequest true "refreshToken request dto"
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.Credential
// @Failure 400 {object} dto.ResponseBadRequestErr "Invalid token"
// @Failure 401 {object} dto.ResponseUnauthorizedErr "Invalid token"
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Router /v1/auth/refreshToken [post]
func (h *handlerImpl) RefreshToken(c router.IContext) {
	request := &dto.RefreshTokenRequest{}
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

	response, respErr := h.service.RefreshToken(request)
	if respErr != nil {
		c.JSON(respErr.StatusCode, respErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

// ForgotPassword is a function to send email to reset password when you forgot password
// @Summary Forgot Password
// @Description Return isSuccess
// @Param request body dto.ForgotPasswordRequest true "forgotPassword request dto"
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.ForgotPasswordResponse
// @Failure 400 {object} dto.ResponseBadRequestErr "Invalid email"
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Router /v1/auth/forgot-password [post]
func (h *handlerImpl) ForgotPassword(c router.IContext) {
	request := &dto.ForgotPasswordRequest{}
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

	response, respErr := h.service.ForgotPassword(request)
	if respErr != nil {
		c.JSON(respErr.StatusCode, respErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

// ResetPassword is a function to reset password
// @Summary Reset Password
// @Description Return isSuccess
// @Param request body dto.ResetPasswordRequest true "resetPassword request dto"
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.ResetPasswordResponse
// @Failure 400 {object} dto.ResponseBadRequestErr "Forbidden the same password"
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Router /v1/auth/reset-password [put]
func (h *handlerImpl) ResetPassword(c router.IContext) {
	request := &dto.ResetPasswordRequest{}
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

	response, errResp := h.service.ResetPassword(request)
	if errResp != nil {
		c.JSON(errResp.StatusCode, errResp)
		return
	}

	c.JSON(http.StatusOK, response)
}
