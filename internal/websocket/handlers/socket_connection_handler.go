package handlers

import (
	"concurrent-image-fetcher/internal/container"
	wsRequest "concurrent-image-fetcher/internal/requests/ws"
	"concurrent-image-fetcher/internal/services"
	wsPkg "concurrent-image-fetcher/internal/websocket/rooms"
	"fmt"
	"github.com/gorilla/websocket"
)

type (
	SocketConnectionHandlerInterface interface {
		HandleWebSocketConnection(*websocket.Conn)
	}
	SocketConnectionHandler struct {
		container       *container.Container
		downloadService services.DownloadServiceInterface
	}
)

// HandleWebSocketConnection listens for messages from the WebSocket connection.
func (handler *SocketConnectionHandler) HandleWebSocketConnection(ws *websocket.Conn) {
	roomId := ws.RemoteAddr().String()
	room := wsPkg.NewRoom(handler.container)
	go room.Listen(ws, roomId)
	for {
		var request wsRequest.DownloadImageRequest
		err := ws.ReadJSON(&request)
		if err != nil {
			fmt.Println("Error al leer desde WebSocket:", err)
			break
		}

		if request.Command == "start_download" {
			go handler.downloadService.ProcessImages(request.Data, roomId)
		}
	}
	room.StopListen()
}

func NewSocketConnectionHandler(container *container.Container) SocketConnectionHandlerInterface {
	return &SocketConnectionHandler{
		container:       container,
		downloadService: services.NewDownloadService(container),
	}
}
