package image

import (
	"net/http"
	"strings"

	"github.com/isd-sgcu/johnjud-gateway/src/app/constant"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	"github.com/isd-sgcu/johnjud-gateway/src/app/router"
	"github.com/isd-sgcu/johnjud-gateway/src/app/validator"
	imageSvc "github.com/isd-sgcu/johnjud-gateway/src/pkg/service/image"
)

type Handler struct {
	service  imageSvc.Service
	validate validator.IDtoValidator
}

func NewHandler(service imageSvc.Service, validate validator.IDtoValidator) *Handler {
	return &Handler{service, validate}
}

// Upload is a function for uploading image to bucket
// @Summary Upload image
// @Description Returns the data of image. If updating pet, add petId. If creating pet, petId is not specified, but keep the imageId.
// @Param image body dto.UploadImageRequest true "upload image request dto"
// @Tags image
// @Accept json
// @Produce json
// @Success 201 {object} dto.ImageResponse
// @Failure 400 {object} dto.ResponseBadRequestErr "Invalid request body"
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Router /v1/images/ [post]
func (h *Handler) Upload(c *router.FiberCtx) {
	request := &dto.UploadImageRequest{}
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

	response, respErr := h.service.Upload(request)
	if respErr != nil {
		c.JSON(respErr.StatusCode, respErr)
		return
	}

	c.JSON(http.StatusCreated, response)
	return
}

// Delete is a function for deleting image from bucket
// @Summary Delete image
// @Description Returns status of deleting image
// @Param id path string true "image id"
// @Tags image
// @Accept json
// @Produce json
// @Success 200 {object} dto.DeleteResponse
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Router /v1/images/{id} [delete]
func (h *Handler) Delete(c *router.FiberCtx) {
	id, err := c.ID()
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

	c.JSON(http.StatusOK, res.Success)
	return
}
