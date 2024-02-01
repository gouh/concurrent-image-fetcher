package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

// UpgradeHandler upgrades the HTTP server connection to the WebSocket protocol.
func UpgradeHandler(wsConnHandler func(*websocket.Conn)) gin.HandlerFunc {
	return func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
			return
		}
		defer ws.Close()

		wsConnHandler(ws)
	}
}
