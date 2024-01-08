package auth

import (
	"net/http"
	"strings"

	"github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	"github.com/isd-sgcu/johnjud-gateway/src/app/router"
	"github.com/isd-sgcu/johnjud-gateway/src/app/validator"
	likeConst "github.com/isd-sgcu/johnjud-gateway/src/constant/like"
	likeSvc "github.com/isd-sgcu/johnjud-gateway/src/pkg/service/like"
)

type Handler struct {
	service  likeSvc.Service
	validate validator.IDtoValidator
}

func NewHandler(service likeSvc.Service, validate validator.IDtoValidator) *Handler {
	return &Handler{service, validate}
}

// FindByUserId is a function that return all petID and userID that user liked.
// @Summary find likes by user id
// @Description Return dto.ResponseSuccess
// @Param id path string true "user id"
// @Tags like
// @Accept json
// @Produce json
// @Success 200 {object} dto.ResponseSuccess
// @Failure 404 {object} dto.ResponseNotfoundErr "user not found"
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Router /v1/likes/ [get]
func (h *Handler) FindByUserId(c router.IContext) {
	id, err := c.Param("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseErr{
			StatusCode: http.StatusInternalServerError,
			Message:    constant.InvalidIDMessage,
			Data:       nil,
		})
		return
	}

	response, respErr := h.service.FindByUserId(id)
	if respErr != nil {
		c.JSON(respErr.StatusCode, respErr)
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess{
		StatusCode: http.StatusOK,
		Message:    likeConst.FindLikeSuccessMessage,
		Data:       response,
	})
	return
}

// Create is a function for creating a `like` for a pet that a user is interested
// @Summary create like
// @Description Return dto.ResponseSuccess
// @Param create body dto.CreateLikeRequest true "create like request"
// @Tags like
// @Accept json
// @Produce json
// @Success 200 {object} dto.ResponseSuccess
// @Failure 400 {object} dto.ResponseBadRequestErr "Invalid request body"
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Router /v1/likes/ [post]
func (h *Handler) Create(c router.IContext) {
	request := &dto.CreateLikeRequest{}
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

	response, respErr := h.service.Create(request)
	if respErr != nil {
		c.JSON(respErr.StatusCode, respErr)
		return
	}

	c.JSON(http.StatusCreated, dto.ResponseSuccess{
		StatusCode: http.StatusCreated,
		Message:    likeConst.CreateLikeSuccessMessage,
		Data:       response,
	})
	return
}

// Create is a function for delete like in database
// @Summary delete like
// @Description Return dto.ResponseSuccess if like is successfully deleted
// @Param id path string true "user id"
// @Tags like
// @Accept json
// @Produce json
// @Success 200 {object} dto.ResponseSuccess
// @Failure 404 {object} dto.ResponseNotfoundErr "like not found"
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Router /v1/likes/ [delete]
func (h *Handler) Delete(c router.IContext) {
	id, err := c.Param("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseErr{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
		return
	}

	res, errRes := h.service.Delete(id)
	if errRes != nil {
		c.JSON(errRes.StatusCode, errRes)
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess{
		StatusCode: http.StatusOK,
		Message:    likeConst.DelteLikeSuccessMessage,
		Data:       res,
	})
	return
}
