package auth

import (
	"net/http"
	"strings"

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
	Update(string, *dto.PetDto) (*proto.Pet, *dto.ResponseErr)
	ChangeView(string) (*proto.Pet, *dto.ResponseErr)
	Delete(string) (bool, *dto.ResponseErr)
}

type ImageService interface {
	FindByPetId(string) (*image_proto.Image, *dto.ResponseErr)
}

func NewHandler(service Service, imageService ImageService, validate *validator.DtoValidator) *Handler {
	return &Handler{service, imageService, validate}
}

func (h *Handler) FindAll(c *router.FiberCtx) {

}

func (h *Handler) FindOne(c *router.FiberCtx) {
	request := &dto.FindOnePetDto{}
	err := c.Bind(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseErr{
			StatusCode: http.StatusBadRequest,
			Message:    "Binding Request Error due to:" + err.Error(),
			Data:       nil,
		})
	}

	if err := h.validate.Validate(request); err != nil {
		var errorMessage []string
		for _, reqErr := range err {
			errorMessage = append(errorMessage, reqErr.Message)
		}
		c.JSON(http.StatusBadRequest, dto.ResponseErr{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Request Body due to:" + strings.Join(errorMessage, ", "),
			Data:       nil,
		})
		return
	}

	response, respErr := h.service.FindOne(request.Id)
	if respErr != nil {
		c.JSON(respErr.StatusCode, respErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) Create(c *router.FiberCtx) {

}

func (h *Handler) Update(c *router.FiberCtx) {

}

func (h *Handler) ChangeView(c *router.FiberCtx) {

}

func (h *Handler) Delete(c *router.FiberCtx) {

}
