package dto

import "net/http"

type ResponseErr struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

type BadReqErrResponse struct {
	Message     string      `json:"message"`
	FailedField string      `json:"failed_field"`
	Value       interface{} `json:"value"`
}

// For docs

type Credential struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL3BiZX..."`
	RefreshToken string `json:"refresh_token" example:"e7e84d54-7518-4..."`
	ExpiresIn    int    `json:"expires_in" example:"3600"`
}

type ResponseBadRequestErr struct {
	StatusCode int                 `json:"status_code" example:"400"`
	Message    string              `json:"message" example:"Invalid request body"`
	Data       []BadReqErrResponse `json:"data"`
}

type ResponseUnauthorizedErr struct {
	StatusCode int         `json:"status_code" example:"401"`
	Message    string      `json:"message" example:"Invalid token"`
	Data       interface{} `json:"data"`
}

type ResponseForbiddenErr struct {
	StatusCode int         `json:"status_code" example:"403"`
	Message    string      `json:"message" example:"Insufficiency permission"`
	Data       interface{} `json:"data"`
}

type ResponseNotfoundErr struct {
	StatusCode int         `json:"status_code" example:"404"`
	Message    string      `json:"message" example:"Not found"`
	Data       interface{} `json:"data"`
}

type ResponseConflictErr struct {
	StatusCode int         `json:"status_code" example:"409"`
	Message    string      `json:"message" example:"Conflict"`
	Data       interface{} `json:"data"`
}

type ResponseInternalErr struct {
	StatusCode int         `json:"status_code" example:"500"`
	Message    string      `json:"message" example:"Internal service error"`
	Data       interface{} `json:"data"`
}

type ResponseServiceDownErr struct {
	StatusCode int         `json:"status_code" example:"503"`
	Message    string      `json:"message" example:"Service is down"`
	Data       interface{} `json:"data"`
}

type ResponseGatewayTimeoutErr struct {
	StatusCode int         `json:"status_code" example:"504"`
	Message    string      `json:"message" example:"Connection timeout"`
	Data       interface{} `json:"data"`
}

func BadRequestError(message string) *ResponseErr {
	return &ResponseErr{http.StatusBadRequest, message, nil}
}

func UnauthorizedError(message string) *ResponseErr {
	return &ResponseErr{http.StatusUnauthorized, message, nil}
}

func ForbiddenError(message string) *ResponseErr {
	return &ResponseErr{http.StatusForbidden, message, nil}
}

func NotFoundError(message string) *ResponseErr {
	return &ResponseErr{http.StatusNotFound, message, nil}
}

func ConflictError(message string) *ResponseErr {
	return &ResponseErr{http.StatusConflict, message, nil}
}

func InternalServerError(message string) *ResponseErr {
	return &ResponseErr{http.StatusInternalServerError, message, nil}
}

func ServiceUnavailableError(message string) *ResponseErr {
	return &ResponseErr{http.StatusServiceUnavailable, message, nil}
}
