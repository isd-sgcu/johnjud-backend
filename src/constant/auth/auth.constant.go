package auth

var ExcludePath = map[string]struct{}{
	"POST /auth/signup":       {},
	"POST /auth/sigin":        {},
	"POST /auth/signout":      {},
	"POST /auth/verify":       {},
	"POST /auth/refreshToken": {},
	"GET /pet/":               {},
	"GET /adopt/":             {},
}
