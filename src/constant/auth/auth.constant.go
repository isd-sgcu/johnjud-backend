package auth

var ExcludePath = map[string]struct{}{
	"POST /auth/signup":          {},
	"POST /auth/signin":          {},
	"POST /auth/verify":          {},
	"POST /auth/forgot-password": {},
	"PUT /auth/reset-password":   {},
	"GET /user/:id":              {},
	"GET /pets":                  {},
	"GET /pets/:id":              {},
	"GET /adopt":                 {},
}

var AdminPath = map[string]struct{}{
	"DELETE /user/:id":            {},
	"GET /pets/admin":             {},
	"POST /pets":                  {},
	"PUT /pets/:id":               {},
	"PUT /pets/:id/visible":       {},
	"DELETE /pets/:id":            {},
	"POST /images/assign/:pet_id": {},
	"DELETE /images/:id":          {},
	"POST /images/":               {},
}

var VersionList = map[string]struct{}{
	"v1": {},
}
