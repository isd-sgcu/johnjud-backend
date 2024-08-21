package image

import (
	"net/http"

	"github.com/isd-sgcu/johnjud-gateway/constant"
	"github.com/isd-sgcu/johnjud-gateway/internal/dto"
	"github.com/isd-sgcu/johnjud-gateway/internal/router"
	"github.com/isd-sgcu/johnjud-gateway/internal/validator"
	"github.com/rs/zerolog/log"
)

type handlerImpl struct {
	service     Service
	validate    validator.IDtoValidator
	maxFileSize int64
}

func NewHandler(service Service, validate validator.IDtoValidator, maxFileSize int64) *handlerImpl {
	return &handlerImpl{
		service:     service,
		validate:    validate,
		maxFileSize: int64(maxFileSize * 1024 * 1024),
	}
}

// Upload is a function for uploading image to bucket
// @Summary Upload image
// @Description Returns the data of image. If updating pet, add petId. If creating pet, petId is not specified, but keep the imageId.
// @Param image body dto.UploadImageRequest true "upload image request dto"
// @Tags image
// @Accept multipart/form-data
// @Produce json
// @Success 201 {object} dto.ImageResponse
// @Failure 400 {object} dto.ResponseBadRequestErr "Invalid request body"
// @Failure 500 {object} dto.ResponseInternalErr "Internal service error"
// @Failure 503 {object} dto.ResponseServiceDownErr "Service is down"
// @Router /v1/images [post]
func (h *handlerImpl) Upload(c *router.FiberCtx) {
	petId := c.GetFormData("pet_id")
	file, err := c.File("file", constant.AllowContentType, h.maxFileSize)
	if err != nil {
		log.Error().
			Err(err).
			Str("service", "image").
			Str("module", "upload").
			Msg("Invalid content")
		c.JSON(http.StatusBadRequest, dto.ResponseErr{
			StatusCode: http.StatusBadRequest,
			Message:    constant.InvalidContentMessage,
			Data:       nil,
		})
		return
	}

	request := &dto.UploadImageRequest{
		Filename: file.Filename,
		File:     file.Data,
		PetId:    petId,
	}

	response, respErr := h.service.Upload(request)
	if respErr != nil {
		c.JSON(respErr.StatusCode, respErr)
		return
	}

	c.JSON(http.StatusCreated, response)
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
func (h *handlerImpl) Delete(c *router.FiberCtx) {
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
}
