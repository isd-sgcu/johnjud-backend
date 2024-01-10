package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	"github.com/isd-sgcu/johnjud-gateway/src/app/router"
	"github.com/isd-sgcu/johnjud-gateway/src/app/validator"
	imageConst "github.com/isd-sgcu/johnjud-gateway/src/constant/image"
	imageSvc "github.com/isd-sgcu/johnjud-gateway/src/pkg/service/image"
)

type Handler struct {
	service  imageSvc.Service
	validate validator.IDtoValidator
}

func NewHandler(service imageSvc.Service, validate validator.IDtoValidator) *Handler {
	return &Handler{service, validate}
}

func (h *Handler) FindByPetId(c *router.FiberCtx) {
	id, err := c.ID()
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseErr{
			StatusCode: http.StatusInternalServerError,
			Message:    constant.InvalidIDMessage,
			Data:       nil,
		})
		return
	}

	response, respErr := h.service.FindByPetId(id)
	if respErr != nil {
		c.JSON(respErr.StatusCode, respErr)
		return
	}

	c.JSON(http.StatusOK, dto.ResponseSuccess{
		StatusCode: http.StatusOK,
		Message:    imageConst.FindImageSuccessMessage,
		Data:       response,
	})
	return
}

func (h *Handler) Upload(c *router.FiberCtx) {
	request := &dto.UploadImageRequest{}
	err := c.Bind(request)
	fmt.Println("request: ", request)
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

	response, respErr := h.service.Upload(request)
	if respErr != nil {
		c.JSON(respErr.StatusCode, respErr)
		return
	}

	c.JSON(http.StatusCreated, dto.ResponseSuccess{
		StatusCode: http.StatusCreated,
		Message:    imageConst.UploadImageSuccessMessage,
		Data:       response,
	})
	return
}

func (h *Handler) Delete(c *router.FiberCtx) {
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
		Message:    imageConst.DelteImageSuccessMessage,
		Data:       res,
	})
	return
}
