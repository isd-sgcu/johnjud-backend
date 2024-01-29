package router

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/isd-sgcu/johnjud-gateway/src/app/dto"
	"github.com/isd-sgcu/johnjud-gateway/src/app/utils"
)

type IContext interface {
	UserID() string
	Role() string
	Bind(interface{}) error
	JSON(int, interface{})
	ID() (string, error)
	Param(string) (string, error)
	Token() string
	Method() string
	Path() string
	StoreValue(string, string)
	Next() error
	Queries() map[string]string
	File(string, map[string]struct{}, int64) (*dto.DecomposedFile, error)
	GetFormData(string) string
}

type FiberCtx struct {
	*fiber.Ctx
}

func NewFiberCtx(c *fiber.Ctx) *FiberCtx {
	return &FiberCtx{c}
}

func (c *FiberCtx) UserID() string {
	return c.Ctx.Locals("UserId").(string)
}

func (c *FiberCtx) Role() string {
	return c.Ctx.Locals("Role").(string)
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

func (c *FiberCtx) Queries() map[string]string {
	return c.Ctx.Queries()
}

func (c *FiberCtx) File(key string, allowContent map[string]struct{}, maxSize int64) (*dto.DecomposedFile, error) {
	file, err := c.Ctx.FormFile(key)
	if err != nil {
		return nil, err
	}

	if !utils.IsExisted(allowContent, file.Header["Content-Type"][0]) {
		return nil, errors.New("Not allow content")
	}

	if file.Size > maxSize {
		return nil, errors.New(fmt.Sprintf("Max file size is %v", maxSize))
	}
	content, err := file.Open()
	if err != nil {
		return nil, errors.New("Cannot read file")
	}

	defer content.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, content); err != nil {
		return nil, err
	}

	return &dto.DecomposedFile{
		Filename: file.Filename,
		Data:     buf.Bytes(),
	}, nil
}

func (c *FiberCtx) GetFormData(key string) string {
	return c.Ctx.FormValue(key)
}

//func (c *FiberCtx) Next() {
//	err := c.Ctx.Next()
//	fmt.Println(c.Route().Path)
//	fmt.Println("next error:", err)
//}
