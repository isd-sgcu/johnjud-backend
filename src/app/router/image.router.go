package router

import "github.com/gofiber/fiber/v2"

func (r *FiberRouter) PostImage(path string, h func(ctx *FiberCtx)) {
	r.image.Post(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}

func (r *FiberRouter) DeleteImage(path string, h func(ctx *FiberCtx)) {
	r.image.Delete(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}

func (r *FiberRouter) GetImages(path string, h func(ctx *FiberCtx)) {
	r.image.Get(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}
