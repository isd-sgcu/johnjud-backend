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
	Adopt(string, *dto.AdoptByRequest) (bool, *dto.ResponseErr)
}

func NewHandler(service Service, imageService imageSrv.Service, validate validator.IDtoValidator) *Handler {
	return &Handler{service, imageService, validate}
}

// FindAll is a function that return all pets in database
// @Summary find all pets
// @Description Return the data of pets if successfully
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.PetDto
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Router /v1/pets/ [get]
func (h *Handler) FindAll(c router.IContext) {
	response, respErr := h.service.FindAll()
	if respErr != nil {
		c.JSON(respErr.StatusCode, respErr)
		return
	}

	c.JSON(http.StatusOK, response)
	return
}

// FindOne is a function that return all pet in database
// @Summary find one pet
// @Description Return the data of pets if successfully
// @Param id path string true "pet id"
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.PetDto
// @Failure 400 {object} dto.ResponseBadRequestErr "Invalid request body"
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Router /v1/pets/{id} [get]
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

// Create is a function that create pet in database
// @Summary create pet
// @Description Return the data of pet if successfully
// @Param create body dto.CreatePetRequest true "pet dto"
// @Tags auth
// @Accept json
// @Produce json
// @Success 201 {object} dto.PetDto
// @Failure 400 {object} dto.ResponseBadRequestErr "Invalid request body"
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Router /v1/pets/create [post]
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

// Update is a function that update pet in database
// @Summary update pet
// @Description Return the data of pet if successfully
// @Param update body dto.UpdatePetRequest true "update pet dto"
// @Param id path stirng true "pet id"
// @Tags auth
// @Accept json
// @Produce json
// @Success 201 {object} dto.PetDto
// @Failure 400 {object} dto.ResponseBadRequestErr "Invalid request body"
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Router /v1/pets/{id} [put]
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

// Change is a function that change visibility of pet in database
// @Summary change view pet
// @Description Return the status true of pet if successfully else false
// @Param change view body dto.ChangeViewPetRequest true "change view pet dto"
// @Param id string true "pet id"
// @Tags auth
// @Accept json
// @Produce json
// @Success 201 {object} bool
// @Failure 400 {object} dto.ResponseBadRequestErr "Invalid request body"
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Router /v1/pets/ [put]
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

	request := &dto.ChangeViewPetRequest{}

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

// Delete is a function that delete pet in database
// @Summary delete pet
// @Description Return the status true of pet if successfully else false
// @Param id string true "pet id"
// @Tags auth
// @Accept json
// @Produce json
// @Success 201 {object} bool
// @Failure 400 {object} dto.ResponseBadRequestErr "Invalid request body"
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Router /v1/pets/ [delete]
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

// Adopt is a function that handles the adoption of a pet in the database
// @Summary Adopt a pet
// @Description Return true if the pet is successfully adopted
// @Param id path string true "Pet ID"
// @Param user_id body string true "User ID"
// @Param pet_id body string true "Pet ID"
// @Tags pet
// @Accept json
// @Produce json
// @Success 201 {object} bool
// @Failure 400 {object} dto.ResponseBadRequestErr "Invalid request body"
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Router /v1/pets/{id}/adopt [put]
func (h *Handler) Adopt(c router.IContext) {
	petId, err := c.Param("id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, &dto.ResponseErr{
			StatusCode: http.StatusInternalServerError,
			Message:    "Invalid ID",
			Data:       nil,
		})
		return
	}

	request := &dto.AdoptByRequest{}
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

	pet, errRes := h.service.Adopt(petId, request)
	if errRes != nil {
		c.JSON(errRes.StatusCode, errRes)
		return
	}

	c.JSON(http.StatusOK, pet)
	return
}
