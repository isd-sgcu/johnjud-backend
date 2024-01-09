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
	"DELETE /users/:id":     {},
	"POST /pets":            {},
	"PUT /pets/:id":         {},
	"PUT /pets/:id/visible": {},
	"DELETE /pets/:id":      {},
	//need to add image upload, delete, assignpet
}

var VersionList = map[string]struct{}{
	"v1": {},
}
