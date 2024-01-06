package auth

var ExcludePath = map[string]struct{}{
	"POST /auth/signup": {},
	"POST /auth/signin": {},
	"POST /auth/verify": {},
	"GET /pet/":         {},
	"GET /adopt/":       {},
}

var VersionList = map[string]struct{}{
	"v1": {},
}
