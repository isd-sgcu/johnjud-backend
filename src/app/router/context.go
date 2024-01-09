package router

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type IContext interface {
	UserID() string
	Bind(interface{}) error
	JSON(int, interface{})
	ID() (string, error)
	Param(string) (string, error)
	Token() string
	Method() string
	Path() string
	StoreValue(string, string)
	IsAuth() bool
	Next() error
}

type FiberCtx struct {
	*fiber.Ctx
}

func NewFiberCtx(c *fiber.Ctx) *FiberCtx {
	return &FiberCtx{c}
}

func (c *FiberCtx) UserID() string {
	userId := c.Ctx.Locals("UserId")
	if userId != nil {
		return userId.(string)
	}
	return ""
	// return c.Ctx.Locals("UserId").(string)
}

func (c *FiberCtx) Bind(v interface{}) error {
	return c.Ctx.BodyParser(v)
}

func (c *FiberCtx) JSON(statusCode int, v interface{}) {
	c.Ctx.Status(statusCode).JSON(v)
}

func (c *FiberCtx) ID() (id string, err error) {
	id = c.Params("id")

	_, err = uuid.Parse(id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (c *FiberCtx) Param(key string) (value string, err error) {
	value = c.Params(key)

	if key == "id" {
		_, err = uuid.Parse(value)
		if err != nil {
			return "", err
		}
	}

	return value, nil
}

func (c *FiberCtx) Token() string {
	raw := c.Ctx.Get(fiber.HeaderAuthorization, "")
	parts := strings.Split(raw, " ")

	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}

func (c *FiberCtx) Method() string {
	return c.Ctx.Method()
}

func (c *FiberCtx) Path() string {
	return c.Ctx.Path()
}

func (c *FiberCtx) StoreValue(k string, v string) {
	c.Locals(k, v)
}

func (c *FiberCtx) IsAuth() bool {
	raw := c.Ctx.Get(fiber.HeaderAuthorization, "")
	parts := strings.Split(raw, " ")

	if len(parts) > 1 {
		return true
	}
	return false
}

//func (c *FiberCtx) Next() {
//	err := c.Ctx.Next()
//	fmt.Println(c.Route().Path)
//	fmt.Println("next error:", err)
//}
