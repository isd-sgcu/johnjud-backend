package router

import "github.com/gofiber/fiber/v2"

func (r *FiberRouter) GetUser(path string, h func(ctx IContext)) {
	r.user.Get(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}

func (r *FiberRouter) PutUser(path string, h func(ctx IContext)) {
	r.user.Put(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}

func (r *FiberRouter) DeleteUser(path string, h func(ctx IContext)) {
	r.user.Delete(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}
