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
	adopt fiber.Router
}

type IGuard interface {
	Use(*FiberCtx)
}

func NewFiberRouter(authGuard IGuard, conf config.App) *FiberRouter {
	r := fiber.New(fiber.Config{
		StrictRouting: true,
		AppName:       "JohnJud API",
	})

	r.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	if conf.Debug {
		r.Use(logger.New(logger.Config{Next: func(c *fiber.Ctx) bool {
			return c.Path() == "/"
		}}))
		r.Get("/docs/*", swagger.HandlerDefault)
	}

	auth := GroupWithAuthMiddleware(r, "/auth", authGuard.Use)
	user := GroupWithAuthMiddleware(r, "/user", authGuard.Use)
	pet := GroupWithAuthMiddleware(r, "/pet", authGuard.Use)
	image := GroupWithAuthMiddleware(r, "/image", authGuard.Use)
	like := GroupWithAuthMiddleware(r, "/like", authGuard.Use)
	adopt := GroupWithAuthMiddleware(r, "/adopt", authGuard.Use)

	return &FiberRouter{r, user, auth, pet, image, adopt, like}
}

func GroupWithAuthMiddleware(r *fiber.App, path string, middleware func(ctx *FiberCtx)) fiber.Router {
	return r.Group(path, func(c *fiber.Ctx) error {
		middleware(NewFiberCtx(c))
		return nil
	})
}
