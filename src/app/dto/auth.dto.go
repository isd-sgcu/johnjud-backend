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

type Signup struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,gte=6,lte=30"`
	Firstname string `json:"firstname" validate:"required"`
	Lastname  string `json:"lastname" validate:"required"`
}

type Signin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=6,lte=30"`
}
