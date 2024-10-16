package api

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
)

func checkAdmin(ctx *gin.Context) {
	session := sessions.Default(ctx)
	role := session.Get("role")
	if role != "admin" {
		log.Panicln("你无权限操作")
	}
}
