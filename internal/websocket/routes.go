package websocket

import (
	"concurrent-image-fetcher/internal/container"
	"concurrent-image-fetcher/internal/websocket/handlers"
	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up the route for handling WebSocket connections.
func SetupRoutes(router *gin.Engine, container *container.Container) {
	socketConnectionHandler := handlers.NewSocketConnectionHandler(container)
	router.GET("/ws", handlers.UpgradeHandler(socketConnectionHandler.HandleWebSocketConnection))
}
