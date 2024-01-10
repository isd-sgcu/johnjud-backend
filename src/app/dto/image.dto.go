package dto

type ImageResponse struct {
	Id        string `json:"id"`
	Url       string `json:"url"`
	ObjectKey string `json:"object_key"`
}

type UploadImageRequest struct {
	Filename string `json:"filename" validate:"required"`
	Data     []byte `json:"data" validate:"required"`
	PetId    string `json:"pet_id" validate:"required"`
}

type DeleteImageResponse struct {
	Success bool `json:"success"`
}
