package router

import (
	"dst-admin-go/api"
	"github.com/gin-gonic/gin"
)

func initLobbyServer(router *gin.RouterGroup) {

	lobbyServerApi := api.LobbyServerApi{}
	router.GET("/lobby/server/query", lobbyServerApi.QueryLobbyServerList)
	router.GET("/lobby/server/query/detail", lobbyServerApi.QueryLobbyServerDetail)

}
