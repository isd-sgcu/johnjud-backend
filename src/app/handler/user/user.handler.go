package user

import (
	"net/http"

	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	"github.com/isd-sgcu/johnjud-gateway/src/app/router"
	"github.com/isd-sgcu/johnjud-gateway/src/app/validator"
	"github.com/isd-sgcu/johnjud-gateway/src/pkg/service/user"
)

type Handler struct {
	service  user.Service
	validate validator.IDtoValidator
}

func NewHandler(service user.Service, validate validator.IDtoValidator) *Handler {
	return &Handler{service, validate}
}

// FindOne is a function that returns a user by id from database
// @Summary finds one user
// @Description Returns the data of user if successful
// @Param id path string true "user id"
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} dto.User
// @Failure 400 {object} dto.ResponseBadRequestErr "Invalid request body"
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Router /v1/users/{id} [get]
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

	usrDto := dto.UpdateUserRequest{}

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
