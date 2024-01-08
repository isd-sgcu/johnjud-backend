package auth

var ExcludePath = map[string]struct{}{
	"POST /auth/signup": {},
	"POST /auth/signin": {},
	"POST /auth/verify": {},
	"GET /pet/":         {},
	"GET /adopt/":       {},
}

var AdminPath = map[string]struct{}{}

var VersionList = map[string]struct{}{
	"v1": {},
}
