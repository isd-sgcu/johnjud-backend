package dto

type LikeResponse struct {
	UserID string `json:"user_id"`
	PetID  string `json:"pet_id"`
}

type FindLikeRequest struct {
	UserID string `json:"user_id" validate:"required"`
}

type CreateLikeRequest struct {
	UserID string `json:"user_id" validate:"required"`
	PetID  string `json:"pet_id" validate:"required"`
}

type DeleteLikeRequest struct {
	Id string `json:"id" validate:"required"`
}

type DeleteLikeResponse struct {
	Success bool `json:"success"`
}
