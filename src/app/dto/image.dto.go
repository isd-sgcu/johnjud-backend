package dto

type DecomposedFile struct {
	Filename string
	Data     []byte
}

type ImageResponse struct {
	Id        string `json:"id"`
	Url       string `json:"url"`
	ObjectKey string `json:"object_key"`
}

type UploadImageRequest struct {
	Filename string `json:"filename" validate:"required"`
	File     []byte `json:"file" validate:"required"`
	PetId    string `json:"pet_id"`
}

type DeleteImageResponse struct {
	Success bool `json:"success"`
}

type AssignPetRequest struct {
	Ids   []string `json:"ids" validate:"required"`
	PetId string   `json:"pet_id" validate:"required"`
}

type AssignPetResponse struct {
	Success bool `json:"success"`
}
