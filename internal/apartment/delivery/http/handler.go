package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sschiz/apartament/internal/apartment"
	"github.com/sschiz/apartament/models"
	"net/http"
	"strconv"
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
	if err := c.ShouldBindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var opts []apartment.Option

	limit, ok := c.GetQuery("limit")

	if ok {
		limit, err := strconv.Atoi(limit)

		if err != nil {
			c.JSON(http.StatusBadRequest,
				map[string]interface{}{
					"error": "wrong limit",
				},
			)
			return
		}

		opts = append(opts, apartment.WithLimit(limit))
	}

	offset, ok := c.GetQuery("offset")
	if ok {
		offset, err := strconv.Atoi(offset)

		if err != nil {
			c.JSON(http.StatusBadRequest,
				map[string]interface{}{
					"error": "wrong offset",
				},
			)
			return
		}

		opts = append(opts, apartment.WithOffset(offset))
	}

	orderField, ok := c.GetQuery("order_field")
	if ok {
		opts = append(opts, apartment.WithOrder(orderField))
	}

	apts, err := h.useCase.Get(c.Request.Context(), inp, opts...)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"apartments": apts,
	})
}
