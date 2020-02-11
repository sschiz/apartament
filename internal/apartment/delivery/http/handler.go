package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sschiz/apartament/internal/apartment"
	"github.com/sschiz/apartament/models"
	"net/http"
)

type Handler struct {
	useCase apartment.UseCase
}

func NewHandler(useCase apartment.UseCase) *Handler {
	return &Handler{useCase: useCase}
}

func (h *Handler) Create(c *gin.Context) {
	inp := new(models.Apartment)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := h.useCase.Create(c.Request.Context(), inp); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) Get(c *gin.Context) {
	inp := new(models.Apartment)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var opts []apartment.Option

	limit, ok := c.Get("limit")
	if ok {
		opts = append(opts, apartment.WithLimit(limit.(int)))
	}

	offset, ok := c.Get("offset")
	if ok {
		opts = append(opts, apartment.WithOffset(offset.(int)))
	}

	orderField, ok := c.Get("order_field")
	if ok {
		opts = append(opts, apartment.WithOrder(orderField.(string)))
	}

	if _, err := h.useCase.Get(c.Request.Context(), inp, opts...); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
