package healthcheck

import (
	"net/http"

	"github.com/isd-sgcu/johnjud-gateway/internal/router"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) HealthCheck(c *router.FiberCtx) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"Health": "ok",
	})
	return
}
