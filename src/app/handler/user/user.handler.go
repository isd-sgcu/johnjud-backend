package user

import (
	"net/http"

	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	"github.com/isd-sgcu/johnjud-gateway/src/app/router"
	"github.com/isd-sgcu/johnjud-gateway/src/app/validator"
	user_proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/user/v1"
)

type Handler struct {
	service  Service
	validate *validator.DtoValidator
}

type Service interface {
	FindOne(string) (*user_proto.User, *dto.ResponseErr)
	Update(string, *dto.UpdateUserDto) (*user_proto.User, *dto.ResponseErr)
	Delete(string) (bool, *dto.ResponseErr)
}

func NewHandler(service Service, validate *validator.DtoValidator) *Handler {
	return &Handler{service, validate}
}

func (h *Handler) FindOne(c *router.FiberCtx) {
	id, err := c.ID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseErr{
			StatusCode: http.StatusInternalServerError,
			Message:    "Invalid ID",
			Data:       nil,
		})
		return
	}

	user, errRes := h.service.FindOne(id)
	if errRes != nil {
		c.JSON(errRes.StatusCode, errRes)
		return
	}

	c.JSON(http.StatusCreated, user)
	return
}

func (h *Handler) Update(c *router.FiberCtx) {
	usrId := c.UserID()

	usrDto := dto.UpdateUserDto{}

	err := c.Bind(&usrDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	user, errRes := h.service.Update(usrId, &usrDto)
	if errRes != nil {
		c.JSON(errRes.StatusCode, errRes)
		return
	}

	c.JSON(http.StatusOK, user)
	return
}

func (h *Handler) Delete(c *router.FiberCtx) {
	id, err := c.ID()
	if err != nil {
		c.JSON(http.StatusBadRequest, &dto.ResponseErr{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
		return
	}

	user, errRes := h.service.Delete(id)
	if errRes != nil {
		c.JSON(errRes.StatusCode, errRes)
		return
	}

	c.JSON(http.StatusOK, user)
	return
}
