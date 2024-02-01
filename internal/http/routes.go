package http

import (
	"concurrent-image-fetcher/internal/container"
	handlers "concurrent-image-fetcher/internal/http/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, container *container.Container) {
	handler := handlers.NewImageHandler(container)
	router.GET("api/v1/images", handler.GetImages)
	router.GET("api/v1/images/:id", handler.GetImage)
	router.POST("api/v1/images", handler.PostImage)
	router.DELETE("api/v1/images/:id", handler.DeleteImage)
}
