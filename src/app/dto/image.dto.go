package dto

type UploadImageRequest struct {
	Filename string `json:"filename" validate:"required"`
	Data     []byte `json:"data" validate:"required"`
	PetId    string `json:"pet_id" validate:"required"`
}

type DeleteImageResponse struct {
	Success bool `json:"success"`
}
