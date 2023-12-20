package user

import (
	"net/http"

	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	"github.com/isd-sgcu/johnjud-gateway/src/app/router"
	"github.com/isd-sgcu/johnjud-gateway/src/app/validator"
	user_proto "github.com/isd-sgcu/johnjud-go-proto/johnjud/auth/user/v1"
)

type Handler struct {
	service  IService
	validate *validator.DtoValidator
}

type IService interface {
	FindOne(string) (*user_proto.User, *dto.ResponseErr)
	Create(*dto.UserDto) (*user_proto.User, *dto.ResponseErr)
	Update(string, *dto.UpdateUserDto) (*user_proto.User, *dto.ResponseErr)
	Delete(string) (bool, *dto.ResponseErr)
}

func NewHandler(service IService, validate *validator.DtoValidator) *Handler {
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

func (h *Handler) Create(c *router.FiberCtx) {
	usrDto := dto.UserDto{}

	err := c.Bind(&usrDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if errors := h.validate.Validate(usrDto); errors != nil {
		c.JSON(http.StatusBadRequest, &dto.ResponseErr{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid body request",
			Data:       errors,
		})
		return
	}

	user, errRes := h.service.Create(&usrDto)
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
