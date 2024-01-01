package dto

type TokenPayloadAuth struct {
	UserId string `json:"user_id"`
}

type Validate struct {
	Token string `json:"token" validate:"jwt"`
}

type RedeemNewToken struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type SignupRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,gte=6,lte=30"`
	Firstname string `json:"firstname" validate:"required"`
	Lastname  string `json:"lastname" validate:"required"`
}

type SignupResponse struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=6,lte=30"`
}

type SignOutResponse struct {
	IsSuccess bool `json:"is_success"`
}
