package middleware

import (
	"github.com/gin-contrib/sessions"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	whitelist = []string{"/api/login", "/api/logout", "/ws", "/api/bootstrap", "/api/init", "/api/install/steamcmd", "/api/kv"}
)

// 拦截检查是否安装 dst steam cmd

func apiFilter(s []string, str string) bool {
	//开放不是 /api 开头接口
	if !strings.Contains(str, "/api") {
		return true
	}
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if apiFilter(whitelist, path) {
			c.Next()
			return
		} else {
			session := sessions.Default(c)
			username := session.Get("username")
			log.Println("auth:", username)
			if username == nil || username == "" {
				c.AbortWithStatus(http.StatusUnauthorized)
			} else {
				c.Next()
			}
		}
	}
}

func SseHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Next()
	}
}
