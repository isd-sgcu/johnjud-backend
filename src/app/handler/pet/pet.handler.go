package pet

import (
	"net/http"
	"strings"

	"github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	imageSvc "github.com/isd-sgcu/johnjud-gateway/src/app/handler/image"
	"github.com/isd-sgcu/johnjud-gateway/src/app/router"
	"github.com/isd-sgcu/johnjud-gateway/src/app/validator"
	petconst "github.com/isd-sgcu/johnjud-gateway/src/constant/pet"
)

type Handler struct {
	service      Service
	imageService imageSvc.Service
	validate     validator.IDtoValidator
}

type Service interface {
	FindAll() ([]*dto.PetResponse, *dto.ResponseErr)
	FindOne(string) (*dto.PetResponse, *dto.ResponseErr)
	Create(*dto.CreatePetRequest) (*dto.PetResponse, *dto.ResponseErr)
	Update(string, *dto.UpdatePetRequest) (*dto.PetResponse, *dto.ResponseErr)
	ChangeView(string, *dto.ChangeViewPetRequest) (*dto.ChangeViewPetResponse, *dto.ResponseErr)
	Delete(string) (*dto.DeleteResponse, *dto.ResponseErr)
}

func NewHandler(service Service, imageService imageSvc.Service, validate validator.IDtoValidator) *Handler {
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

	c.JSON(http.StatusOK, dto.ResponseSuccess{
		StatusCode: http.StatusOK,
		Message:    petconst.FindAllPetSuccessMessage,
		Data:       response,
	})
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
		c.JSON(http.StatusInternalServerError, dto.ResponseErr{
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

	c.JSON(http.StatusOK, dto.ResponseSuccess{
		StatusCode: http.StatusOK,
		Message:    petconst.FindOnePetSuccessMessage,
		Data:       response,
	})
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

	c.JSON(http.StatusCreated, dto.ResponseSuccess{
		StatusCode: http.StatusCreated,
		Message:    petconst.CreatePetSuccessMessage,
		Data:       response,
	})
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

	c.JSON(http.StatusOK, dto.ResponseSuccess{
		StatusCode: http.StatusOK,
		Message:    petconst.UpdatePetSuccessMessage,
		Data:       pet,
	})
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

	c.JSON(http.StatusOK, dto.ResponseSuccess{
		StatusCode: http.StatusOK,
		Message:    petconst.ChangeViewPetSuccessMessage,
		Data:       res,
	})
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

	c.JSON(http.StatusOK, dto.ResponseSuccess{
		StatusCode: http.StatusOK,
		Message:    petconst.DeletePetSuccessMessage,
		Data:       res,
	})
	return
}
