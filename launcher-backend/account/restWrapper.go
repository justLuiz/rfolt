package account

import (
	"github.com/rocwong/neko"
	"rfolt/launcher-backend/config"
)

func Check(ctx *neko.Context) {
	//Получим логин и пароль с запроса
	login := ctx.Params.ByPost("login")
	password := ctx.Params.ByPost("password")

	//Проверим на корректность
	if (len(login) < 3 || len(login) > 13) || (len(password) < 3 || len(password) > 13) {
		ctx.Json(neko.JSON{"code": -1, "error": "Неверный логин или пароль"})
		return
	}

	//Получим статус аккаунта
	result := loginServerWork(config.GetNetwork(), login, password, StopOnCheck)
	if result.NetworkState == -1 {
		ctx.Json(neko.JSON{"code": -2, "error": "Ошибка при подключении к серверу"})
	} else if result.NetworkState == -2 {
		ctx.Json(neko.JSON{"code": -3, "error": "Ошибка при получении данных от сервера"})
	} else {
		ctx.Json(neko.JSON{"code": 0, "accountStatus": result.AccountStatus})
	}
}

func StartGame(ctx *neko.Context) {
	//Получим логин и пароль с запроса
	login := ctx.Params.ByPost("login")
	password := ctx.Params.ByPost("password")

	//Проверим на корректность
	if (len(login) < 3 || len(login) > 13) || (len(password) < 3 || len(password) > 13) {
		ctx.Json(neko.JSON{"code": -1, "error": "Неверный логин или пароль"})
		return
	}

	//Получим статус аккаунта
	result := loginServerWork(config.GetNetwork(), login, password, -1)
	if result.NetworkState == -1 {
		ctx.Json(neko.JSON{"code": -2, "error": "Ошибка при подключении к серверу"})
	} else if result.NetworkState == -2 {
		ctx.Json(neko.JSON{"code": -3, "error": "Ошибка при получении данных от сервера"})
	} else {
		ctx.Json(neko.JSON{"code": 0, "accountStatus": result.AccountStatus, "pid": result.BinPid})
	}
}
