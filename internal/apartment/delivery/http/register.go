package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sschiz/apartament/internal/apartment"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, uc apartment.UseCase) {
	h := NewHandler(uc)

	apartments := router.Group("/apartments")
	{
		apartments.POST("", h.Create)
		apartments.GET("", h.Get)
	}
}
