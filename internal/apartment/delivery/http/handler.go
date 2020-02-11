package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sschiz/apartament/internal/apartment"
	"github.com/sschiz/apartament/models"
	"log"
	"net/http"
	"os"
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
	inp := c.MustGet("apt").(*models.Apartment)

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

func Validation() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "GET" {
			c.Next()
			return
		}

		inp := c.MustGet("apt").(*models.Apartment)

		if inp.Rent < 1 || inp.Rooms < 1 || inp.Area < 1 || inp.Floor < 1 || inp.House == nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		house := inp.House
		if house.Floors < 1 {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if ac := inp.ApartmentComplex; ac != nil {
			if len(ac.Name) == 0 || ac.Apartments[0] < 1 || ac.Apartments[1] < 1 {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
		}

		c.Next()
	}
}

func Logger() gin.HandlerFunc {
	logger := log.New(os.Stdout, "[LOGGER] ", log.LstdFlags)
	return func(c *gin.Context) {
		apt := new(models.Apartment)
		if err := c.ShouldBindJSON(apt); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		logger.Printf("Got Apartment: %#v", apt)
		logger.Printf("House: %#v", apt.House)
		logger.Printf("Apartment complex: %#v", apt.ApartmentComplex)

		c.Set("apt", apt)

		c.Next()
	}
}
