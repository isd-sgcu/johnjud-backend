package auth

import (
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	"github.com/isd-sgcu/johnjud-gateway/src/app/router"
	"github.com/isd-sgcu/johnjud-gateway/src/app/validator"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/file/image/v1"
)

type Handler struct {
	service  Service
	validate *validator.DtoValidator
}

type Service interface {
	FindByPetId(string) ([]*proto.Image, *dto.ResponseErr)
	Upload(*dto.ImageDto) (*proto.Image, *dto.ResponseErr)
	Delete(string) (bool, *dto.ResponseErr)
}

func NewHandler(service Service, validate *validator.DtoValidator) *Handler {
	return &Handler{service, validate}
}

func (h *Handler) FindByPetId(c *router.FiberCtx) {

}

func (h *Handler) Upload(c *router.FiberCtx) {

}

func (h *Handler) Delete(c *router.FiberCtx) {

}
