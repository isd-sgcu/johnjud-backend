package router

import "github.com/gofiber/fiber/v2"

func (r *FiberRouter) GetAdopt(path string, h func(ctx *FiberCtx)) {
	r.adopt.Get(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}

func (r *FiberRouter) PostAdopt(path string, h func(ctx *FiberCtx)) {
	r.adopt.Post(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}

func (r *FiberRouter) DeleteAdopt(path string, h func(ctx *FiberCtx)) {
	r.adopt.Delete(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}
