package rooms

import (
	"concurrent-image-fetcher/internal/container"
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/gorilla/websocket"
)

type (
	RoomInterface interface {
		Listen(*websocket.Conn, string)
		StopListen()
	}
	Room struct {
		Redis  *redis.Client
		Cancel context.CancelFunc
		Ctx    context.Context
	}
)

func (room *Room) IsListen() bool {
	if room.Ctx != nil {
		return room.Ctx.Err() == nil
	}
	return false
}

func (room *Room) Listen(ws *websocket.Conn, roomId string) {
	if !room.IsListen() {
		room.Ctx, room.Cancel = context.WithCancel(context.Background())
		pubSub := room.Redis.Subscribe(room.Ctx, roomId)
		defer pubSub.Close()
		for {
			select {
			case msg := <-pubSub.Channel():
				err := ws.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
				if err != nil {
					fmt.Println(err.Error())
				}
			case <-room.Ctx.Done():
				break
			}
		}
	}
}

func (room *Room) StopListen() {
	if room.IsListen() {
		room.Cancel()
	}
}

func NewRoom(container *container.Container) RoomInterface {
	return &Room{
		Redis: container.Redis,
	}
}
