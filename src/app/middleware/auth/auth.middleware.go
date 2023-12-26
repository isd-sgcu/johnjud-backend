package auth

import (
	"github.com/isd-sgcu/johnjud-gateway/src/app/router"
	"github.com/isd-sgcu/johnjud-gateway/src/config"
	authPkg "github.com/isd-sgcu/johnjud-gateway/src/pkg/service/auth"
)

type Guard struct {
	service  authPkg.Service
	excludes map[string]struct{}
	conf     config.App
}

func NewAuthGuard(s authPkg.Service, e map[string]struct{}, conf config.App) Guard {
	return Guard{
		service:  s,
		excludes: e,
		conf:     conf,
	}
}

func (m *Guard) Use(ctx router.IContext) error {
	err := m.Validate(ctx)
	if err != nil {
		return err
	}

	return ctx.Next()

}

func (m *Guard) Validate(ctx router.IContext) error {
	return nil
}
