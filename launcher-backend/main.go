package main

import (
	"flag"
	"fmt"
	"github.com/rocwong/neko"
	"net/http"
	"rfolt/launcher-backend/account"
	"rfolt/launcher-backend/config"
)

func main() {
	//Распарсим флаги коммандной строки
	debug := flag.Bool("debug", false, "Run in debug mode")
	port := flag.Int("port", 1488, "UI-server port")
	serverId := flag.String("id", "", "RFOLt server id")
	flag.Parse()

	config.SetServerId(*serverId)

	//todo Менеджер обновлений

	//Менеджер игрового клиента

	//Запустим сервер комманд
	server := neko.New()
	if *debug {
		server.Use(neko.Logger())
	}

	//Роуты
	server.GET("/version", func(ctx *neko.Context) { //версия
		ctx.Json(neko.JSON{
			"core": "2018.9",
		}, http.StatusOK)
	})

	server.GET("/title", func(ctx *neko.Context) {
		ctx.Text(config.GetMain().Title, http.StatusOK)
	})

	//Управление аккаунтами
	server.Group("/account", func(router *neko.RouterGroup) {
		router.POST("/check", account.Check)
		router.POST("/startGame", account.StartGame)
		router.POST("/register")
	})

	server.Run(fmt.Sprintf("127.0.0.1:%v", *port))
}
