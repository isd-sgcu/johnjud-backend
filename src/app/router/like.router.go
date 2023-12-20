package router

import "github.com/gofiber/fiber/v2"

func (r *FiberRouter) GetLike(path string, h func(ctx *FiberCtx)) {
	r.like.Get(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}

func (r *FiberRouter) PostLike(path string, h func(ctx *FiberCtx)) {
	r.like.Post(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}

func (r *FiberRouter) DeleteLike(path string, h func(ctx *FiberCtx)) {
	r.like.Delete(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}
