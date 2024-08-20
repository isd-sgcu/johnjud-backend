package auth

import (
	"net/http"

	"github.com/isd-sgcu/johnjud-gateway/config"
	"github.com/isd-sgcu/johnjud-gateway/constant"
	"github.com/isd-sgcu/johnjud-gateway/internal/dto"
	authPkg "github.com/isd-sgcu/johnjud-gateway/internal/pkg/service/auth"
	"github.com/isd-sgcu/johnjud-gateway/internal/router"
	"github.com/isd-sgcu/johnjud-gateway/internal/utils"
	"github.com/isd-sgcu/johnjud-gateway/internal/utils/auth"
)

type Guard struct {
	service     authPkg.Service
	excludes    map[string]struct{}
	adminpath   map[string]struct{}
	conf        config.App
	versionList map[string]struct{}
}

func NewAuthGuard(s authPkg.Service, e map[string]struct{}, a map[string]struct{}, conf config.App, versionList map[string]struct{}) Guard {
	return Guard{
		service:     s,
		excludes:    e,
		adminpath:   a,
		conf:        conf,
		versionList: versionList,
	}
}

func (m *Guard) Use(ctx router.IContext) error {
	method := ctx.Method()
	path := ctx.Path()

	path = utils.TrimInList(path, "/", m.versionList)
	ids := auth.FindIDFromPath(path)
	path = auth.FormatPath(method, path, ids)
	if utils.IsExisted(m.excludes, path) {
		return ctx.Next()
	}

	token := ctx.Token()
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, &dto.ResponseErr{
			StatusCode: http.StatusUnauthorized,
			Message:    "Invalid token",
		})
		return nil
	}

	payload, err := m.service.Validate(token)
	if err != nil {
		ctx.JSON(err.StatusCode, err)
		return nil
	}

	ctx.StoreValue("UserId", payload.UserId)
	ctx.StoreValue("Role", payload.Role)

	if utils.IsExisted(m.adminpath, path) && payload.Role != string(constant.ADMIN) {
		ctx.JSON(http.StatusUnauthorized, dto.ResponseErr{
			StatusCode: http.StatusUnauthorized,
			Message:    "Limited access",
			Data:       nil,
		})
		return nil
	}

	return ctx.Next()
}
