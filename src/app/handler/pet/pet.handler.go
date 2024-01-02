package pet

import (
	"net/http"
	"strings"

	"github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	imageSrv "github.com/isd-sgcu/johnjud-gateway/src/app/handler/image"
	"github.com/isd-sgcu/johnjud-gateway/src/app/router"
	"github.com/isd-sgcu/johnjud-gateway/src/app/validator"
	pet_proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
)

type Handler struct {
	service      Service
	imageService imageSrv.Service
	validate     validator.IDtoValidator
}

type Service interface {
	FindAll() ([]*pet_proto.Pet, *dto.ResponseErr)
	FindOne(string) (*pet_proto.Pet, *dto.ResponseErr)
	Create(*dto.CreatePetRequest) (*pet_proto.Pet, *dto.ResponseErr)
	Update(string, *dto.UpdatePetRequest) (*pet_proto.Pet, *dto.ResponseErr)
	ChangeView(string, *dto.ChangeViewPetRequest) (bool, *dto.ResponseErr)
	Delete(string) (bool, *dto.ResponseErr)
}

func NewHandler(service Service, imageService imageSrv.Service, validate validator.IDtoValidator) *Handler {
	return &Handler{service, imageService, validate}
}

func (h *Handler) FindAll(c router.IContext) {
	response, respErr := h.service.FindAll()
	if respErr != nil {
		c.JSON(respErr.StatusCode, respErr)
		return
	}

	c.JSON(http.StatusOK, response)
	return
}

func (h *Handler) FindOne(c router.IContext) {
	id, err := c.Param("id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, &dto.ResponseErr{
			StatusCode: http.StatusInternalServerError,
			Message:    "Invalid ID",
			Data:       nil,
		})
		return
	}

	response, respErr := h.service.FindOne(id)
	if respErr != nil {
		c.JSON(respErr.StatusCode, respErr)
		return
	}

	c.JSON(http.StatusOK, response)
	return
}

func (h *Handler) Create(c router.IContext) {
	request := &dto.CreatePetRequest{
		Pet: &dto.PetDto{},
	}
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

	c.JSON(http.StatusCreated, response)
	return
}

func (h *Handler) Update(c router.IContext) {
	petId, err := c.Param("id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, &dto.ResponseErr{
			StatusCode: http.StatusInternalServerError,
			Message:    "Invalid ID",
			Data:       nil,
		})
		return
	}

	request := &dto.UpdatePetRequest{
		Pet: &dto.PetDto{},
	}

	err = c.Bind(request)
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

	pet, errRes := h.service.Update(petId, request)
	if errRes != nil {
		c.JSON(errRes.StatusCode, errRes)
		return
	}

	c.JSON(http.StatusOK, pet)
	return
}

func (h *Handler) ChangeView(c router.IContext) {
	id, err := c.Param("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, &dto.ResponseErr{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
		return
	}

	request := &dto.ChangeViewPetRequest{
		Visible: false,
	}

	err = c.Bind(request)
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

	res, errRes := h.service.ChangeView(id, request)
	if errRes != nil {
		c.JSON(errRes.StatusCode, errRes)
		return
	}

	c.JSON(http.StatusOK, res)
	return
}

func (h *Handler) Delete(c router.IContext) {
	id, err := c.Param("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, &dto.ResponseErr{
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

	c.JSON(http.StatusOK, res)
	return
}
