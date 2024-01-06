package router

import "github.com/gofiber/fiber/v2"

func (r *FiberRouter) GetPet(path string, h func(ctx IContext)) {
	r.pet.Get(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}

func (r *FiberRouter) PostPet(path string, h func(ctx IContext)) {
	r.pet.Post(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}

func (r *FiberRouter) PutPet(path string, h func(ctx IContext)) {
	r.pet.Put(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}

func (r *FiberRouter) DeletePet(path string, h func(ctx IContext)) {
	r.pet.Delete(path, func(c *fiber.Ctx) error {
		h(NewFiberCtx(c))
		return nil
	})
}
