package auth

import (
	"github.com/isd-sgcu/johnjud-gateway/src/app/router"
	"github.com/isd-sgcu/johnjud-gateway/src/app/validator"
	imageSvc "github.com/isd-sgcu/johnjud-gateway/src/pkg/service/image"
)

type Handler struct {
	service  imageSvc.Service
	validate *validator.DtoValidator
}

func NewHandler(service imageSvc.Service, validate *validator.DtoValidator) *Handler {
	return &Handler{service, validate}
}

func (h *Handler) FindByPetId(c *router.FiberCtx) {

}

func (h *Handler) Upload(c *router.FiberCtx) {

}

func (h *Handler) Delete(c *router.FiberCtx) {

}
