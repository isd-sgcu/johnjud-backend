package router

import "github.com/gofiber/fiber/v2"

func (r *FiberRouter) PostAuth(path string, h func(ctx *FiberCtx)) {
	r.auth.Post(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}
