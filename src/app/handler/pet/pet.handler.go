package auth

import (
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

}

func (h *Handler) Create(c *router.FiberCtx) {

}

func (h *Handler) Update(c *router.FiberCtx) {

}

func (h *Handler) ChangeView(c *router.FiberCtx) {

}

func (h *Handler) Delete(c *router.FiberCtx) {

}
