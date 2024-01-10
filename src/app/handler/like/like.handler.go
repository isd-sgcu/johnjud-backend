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

func (h *Handler) FindByUserId(c router.IContext) {
	id, err := c.ID()
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

func (h *Handler) Delete(c router.IContext) {
	id, err := c.ID()
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
