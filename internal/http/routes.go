package http

import (
	"concurrent-image-fetcher/internal/container"
	handlers "concurrent-image-fetcher/internal/http/handlers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func SetupRoutes(router *gin.Engine, container *container.Container) {
	router.Use(corsMiddleware())
	handler := handlers.NewImageHandler(container)
	router.Static("/app", "./web")
	router.Static("/public", "./public")
	router.GET("/api/v1/images", handler.GetImages)
	router.GET("/api/v1/images/:id", handler.GetImage)
	router.POST("/api/v1/images", handler.PostImage)
	router.DELETE("/api/v1/images/:id", handler.DeleteImage)
}
