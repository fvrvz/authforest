package main

import (
	"github.com/fvrvz/authforest/config"
	"github.com/fvrvz/authforest/constants"
	"github.com/fvrvz/authforest/db"
	"github.com/fvrvz/authforest/helpers"
	"github.com/fvrvz/authforest/server"
)

func main() {
	config.Init(constants.CONFIG_PATH, constants.ENV_PATH)
	helpers.InitRSAKey(config.GetConfig().OIDC.RSAKeyPath)
	db.Init()
	server.InitServer()
}
