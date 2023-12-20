package router

import "github.com/gofiber/fiber/v2"

func (r *FiberRouter) GetUser(path string, h func(ctx *FiberCtx)) {
	r.user.Get(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}

func (r *FiberRouter) PutUser(path string, h func(ctx *FiberCtx)) {
	r.user.Put(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}
