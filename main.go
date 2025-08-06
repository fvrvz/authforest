package main

import (
	"github.com/fvrvz/auth-service-go/config"
	"github.com/fvrvz/auth-service-go/constants"
	"github.com/fvrvz/auth-service-go/db"
	"github.com/fvrvz/auth-service-go/server"
)

func main() {
	config.Init(constants.CONFIG_PATH, constants.ENV_PATH)
	db.Init()
	server.InitServer()
}
