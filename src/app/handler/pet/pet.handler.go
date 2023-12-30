package pet

import (
	"net/http"
	"strings"

	"github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	"github.com/isd-sgcu/johnjud-gateway/src/app/router"
	"github.com/isd-sgcu/johnjud-gateway/src/app/validator"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/pet/v1"
	image_proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
)

type Handler struct {
	service      Service
	imageService ImageService
	validate     *validator.DtoValidator
}

type Service interface {
	FindAll() ([]*proto.Pet, *dto.ResponseErr)
	FindOne(string) (*proto.Pet, *dto.ResponseErr)
	Create(*dto.PetDto) (*proto.Pet, *dto.ResponseErr)
	Update(string, *dto.UpdatePetDto) (*proto.Pet, *dto.ResponseErr)
	ChangeView(string, bool) (bool, *dto.ResponseErr)
	Delete(string) (bool, *dto.ResponseErr)
}

type ImageService interface {
	FindByPetId(string) (*image_proto.Image, *dto.ResponseErr)
}

func NewHandler(service Service, imageService ImageService, validate *validator.DtoValidator) *Handler {
	return &Handler{service, imageService, validate}
}

func (h *Handler) FindAll(c router.IContext) {
	request := &dto.FindOnePetDto{}
	err := c.Bind(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseErr{
			StatusCode: http.StatusBadRequest,
			Message:    constant.BindingRequestErrorMessage + err.Error(),
			Data:       nil,
		})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (h *Handler) FindOne(c router.IContext) {
	id, err := c.ID()

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
	request := &dto.PetDto{}
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
	petId, err := c.ID()

	if err != nil {
		c.JSON(http.StatusInternalServerError, &dto.ResponseErr{
			StatusCode: http.StatusInternalServerError,
			Message:    "Invalid ID",
			Data:       nil,
		})
		return
	}

	petDto := dto.UpdatePetDto{}

	err = c.Bind(&petDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	pet, errRes := h.service.Update(petId, &petDto)
	if errRes != nil {
		c.JSON(errRes.StatusCode, errRes)
		return
	}

	c.JSON(http.StatusOK, pet)
	return
}

func (h *Handler) ChangeView(c *router.FiberCtx) {

}

func (h *Handler) Delete(c *router.FiberCtx) {
	id, err := c.Param("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, &dto.ResponseErr{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
		return
	}

	pet, errRes := h.service.Delete(id)
	if errRes != nil {
		c.JSON(errRes.StatusCode, errRes)
		return
	}

	c.JSON(http.StatusOK, pet)
	return
}
