package pet

import (
	"net/http"
	"strings"

	"github.com/isd-sgcu/johnjud-backend/constant"
	"github.com/isd-sgcu/johnjud-backend/internal/dto"
	"github.com/isd-sgcu/johnjud-backend/internal/image"
	"github.com/isd-sgcu/johnjud-backend/internal/router"
	"github.com/isd-sgcu/johnjud-backend/internal/validator"
)

type handlerImpl struct {
	service      Service
	imageService image.Service
	validate     validator.IDtoValidator
}

func NewHandler(service Service, imageService image.Service, validate validator.IDtoValidator) *handlerImpl {
	return &handlerImpl{service, imageService, validate}
}

// FindAll is a function that returns all VISIBLE pets in database
// @Summary finds all pets
// @Description Returns the data of pets if successful
// @Tags pet
// @Accept json
// @Produce json
// @Success 200 {object} []dto.PetResponse
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Router /v1/pets/ [get]
func (h *handlerImpl) FindAll(c router.IContext) {
	queries := c.Queries()
	request, err := QueriesToFindAllDto(queries)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	response, respErr := h.service.FindAll(request, false)
	if respErr != nil {
		c.JSON(respErr.StatusCode, respErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

// FindAllAdmin is a function that returns ALL pets in database
// @Summary finds all pets
// @Description Returns the data of pets if successful
// @Tags pet
// @Accept json
// @Produce json
// @Success 200 {object} []dto.PetResponse
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Router /v1/pets/admin [get]
func (h *handlerImpl) FindAllAdmin(c router.IContext) {
	queries := c.Queries()
	request, err := QueriesToFindAllDto(queries)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	response, respErr := h.service.FindAll(request, true)
	if respErr != nil {
		c.JSON(respErr.StatusCode, respErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

// FindOne is a function that returns a pet by id in database
// @Summary finds one pet
// @Description Returns the data of a pet if successful
// @Param id path string true "pet id"
// @Tags pet
// @Accept json
// @Produce json
// @Success 200 {object} dto.PetResponse
// @Failure 400 {object} dto.ResponseBadRequestErr "Invalid request body"
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Router /v1/pets/{id} [get]
func (h *handlerImpl) FindOne(c router.IContext) {
	id, err := c.Param("id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ResponseErr{
			StatusCode: http.StatusInternalServerError,
			Message:    constant.InvalidIDMessage,
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
}

// Create is a function that creates pet in database
// @Summary creates pet
// @Description Returns the data of pet if successful
// @Param create body dto.CreatePetRequest true "pet dto"
// @Tags pet
// @Accept json
// @Produce json
// @Success 201 {object} dto.PetResponse
// @Failure 400 {object} dto.ResponseBadRequestErr "Invalid request body"
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Router /v1/pets/create [post]
func (h *handlerImpl) Create(c router.IContext) {
	request := &dto.CreatePetRequest{}
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
}

// Update is a function that updates pet in database
// @Summary updates pet
// @Description Returns the data of pet if successfully
// @Param update body dto.UpdatePetRequest true "update pet dto"
// @Param id path string true "pet id"
// @Tags pet
// @Accept json
// @Produce json
// @Success 201 {object} dto.PetResponse
// @Failure 400 {object} dto.ResponseBadRequestErr "Invalid request body"
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Router /v1/pets/{id} [put]
func (h *handlerImpl) Update(c router.IContext) {
	petId, err := c.Param("id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, &dto.ResponseErr{
			StatusCode: http.StatusInternalServerError,
			Message:    "Invalid ID",
			Data:       nil,
		})
		return
	}

	request := &dto.UpdatePetRequest{}

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
}

// ChangeView is a function that changes visibility of pet in database
// @Summary changes pet's public visiblility
// @Description Returns successful status if pet's IsVisible is successfully changed
// @Param changeViewDto body dto.ChangeViewPetRequest true "changeView pet dto"
// @Param id path string true "pet id"
// @Tags pet
// @Accept json
// @Produce json
// @Success 201 {object} dto.ChangeViewPetResponse
// @Failure 400 {object} dto.ResponseBadRequestErr "Invalid request body"
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Router /v1/pets/{id}/visible [put]
func (h *handlerImpl) ChangeView(c router.IContext) {
	id, err := c.Param("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseErr{
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
}

// Delete is a function that deletes pet in database
// @Summary deletes pet
// @Description Returns successful status if pet is successfully deleted
// @Param id path string true "pet id"
// @Tags pet
// @Accept json
// @Produce json
// @Success 201 {object} dto.DeleteResponse
// @Failure 400 {object} dto.ResponseBadRequestErr "Invalid request body"
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Router /v1/pets/{id} [delete]
func (h *handlerImpl) Delete(c router.IContext) {
	id, err := c.Param("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ResponseErr{
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
}

// Adopt is a function that handles the adoption of a pet in the database
// @Summary Change a pet's adoptBy status
// @Description Return true if the pet is successfully adopted
// @Param adoptDto body dto.AdoptByRequest true "adopt pet dto"
// @Param id path string true "Pet ID"
// @Tags pet
// @Accept json
// @Produce json
// @Success 201 {object} dto.AdoptByResponse
// @Failure 400 {object} dto.ResponseBadRequestErr "Invalid request body"
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Router /v1/pets/{id}/adopt [put]
func (h *handlerImpl) Adopt(c router.IContext) {
	petId, err := c.Param("id")
	if err != nil {
		c.JSON(http.StatusBadRequest, &dto.ResponseErr{
			StatusCode: http.StatusBadRequest,
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

	res, errRes := h.service.Adopt(petId, request)
	if errRes != nil {
		c.JSON(errRes.StatusCode, errRes)
		return
	}

	c.JSON(http.StatusOK, res)
}
