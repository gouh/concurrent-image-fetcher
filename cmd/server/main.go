package main

import (
	"concurrent-image-fetcher/config"
	"concurrent-image-fetcher/internal/container"
	"concurrent-image-fetcher/internal/http"
	"concurrent-image-fetcher/internal/websocket"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	envFile := flag.String("env-file", ".env", ".env configuration file path")
	configs := config.NewConfig(*envFile)
	router := gin.Default()
	cont := container.NewContainer(configs)
	http.SetupRoutes(router, cont)
	websocket.SetupRoutes(router, cont)
	err := router.Run(":8080")
	if err != nil {
		fmt.Println(err.Error())
	}
}
