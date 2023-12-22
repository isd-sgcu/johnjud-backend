package auth

import (
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	"github.com/isd-sgcu/johnjud-gateway/src/app/router"
	"github.com/isd-sgcu/johnjud-gateway/src/app/validator"
	proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/backend/like/v1"
)

type Handler struct {
	service  Service
	validate *validator.DtoValidator
}

type Service interface {
	FindByUserId(string) ([]*proto.Like, *dto.ResponseErr)
	Create(*dto.LikeDto) (*proto.Like, *dto.ResponseErr)
	Delete(string) (bool, *dto.ResponseErr)
}

func NewHandler(service Service, validate *validator.DtoValidator) *Handler {
	return &Handler{service, validate}
}

func (h *Handler) FindByUserId(c *router.FiberCtx) {

}

func (h *Handler) Create(c *router.FiberCtx) {

}

func (h *Handler) Delete(c *router.FiberCtx) {

}
