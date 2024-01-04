package router

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/isd-sgcu/johnjud-gateway/src/config"
)

type FiberRouter struct {
	*fiber.App
	auth  fiber.Router
	user  fiber.Router
	pet   fiber.Router
	image fiber.Router
	like  fiber.Router
}

type IGuard interface {
	Use(IContext) error
}

func NewAPIv1(r *FiberRouter, conf config.App) *fiber.App {
	if conf.Debug {
		r.Use(logger.New(logger.Config{Next: func(c *fiber.Ctx) bool {
			return c.Path() == "/v1/"
		}}))
		r.Get("/docs/*", swagger.HandlerDefault)
	}

	app := fiber.New(fiber.Config{
		StrictRouting:     true,
		AppName:           "JohnJud API",
		EnablePrintRoutes: conf.Debug,
	})

	app.Mount("/v1", r.App)

	return app
}

func NewFiberRouter(authGuard IGuard, conf config.App) *FiberRouter {
	r := fiber.New(fiber.Config{})

	r.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	auth := GroupWithAuthMiddleware(r, "/auth", authGuard.Use)
	user := GroupWithAuthMiddleware(r, "/user", authGuard.Use)
	pet := GroupWithAuthMiddleware(r, "/pets", authGuard.Use)

	image := GroupWithAuthMiddleware(r, "/image", authGuard.Use)
	like := GroupWithAuthMiddleware(r, "/like", authGuard.Use)

	return &FiberRouter{r, auth, user, pet, image, like}
}

func GroupWithAuthMiddleware(r *fiber.App, path string, middleware func(ctx IContext) error) fiber.Router {
	return r.Group(path, func(c *fiber.Ctx) error {
		return middleware(NewFiberCtx(c))
	})
}
