package dto

type AdoptDto struct {
	UserID string `json:"user_id" validate:"required"`
	PetID  string `json:"pet_id" validate:"required"`
}

type AdoptByRequest struct {
	Adopt AdoptDto `json:"adopt" validate:"required"`
}
