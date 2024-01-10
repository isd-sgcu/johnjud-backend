package auth

var ExcludePath = map[string]struct{}{
	"POST /auth/signup": {},
	"POST /auth/signin": {},
	"POST /auth/verify": {},
	"GET /user/:id":     {},
	"GET /pets":         {},
	"GET /adopt":        {},
}

var AdminPath = map[string]struct{}{
	"DELETE /user/:id":      {},
	"POST /pets":            {},
	"PUT /pets/:id":         {},
	"PUT /pets/:id/visible": {},
	"DELETE /pets/:id":      {},
	//need to add image upload, delete, assignpet
	"POST /images/assign/:pet_id": {},
	"DELETE /images/:id":          {},
	"POST /images/":               {},
	"GET /images/:id":             {},
}

var VersionList = map[string]struct{}{
	"v1": {},
}
