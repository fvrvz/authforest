package server

import (
	"fmt"
	"log"

	"github.com/fvrvz/authforest/config"
)

func InitServer() {
	router := InitRouter()

	port := config.GetConfig().Server.Port

	addr := fmt.Sprintf(":%d", port)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
