package middleware

import (
	"dst-admin-go/internal/service/login"
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	whitelist = []string{"/api/login", "/api/logout", "/ws", "/api/dst-static", "/api/init", "/api/install/steamcmd"}
)

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
func Authentication(loginService *login.LoginService) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if apiFilter(whitelist, path) {
			c.Next()
			return
		}
		session := sessions.Default(c)
		username := session.Get("username")
		log.Println("username:", username)
		if username == nil {
			if loginService.IsWhiteIP(c) {
				loginService.DirectLogin(c)
				log.Println("white ip login", c.ClientIP())
				c.Next()
				return
			}
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			c.Next()
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
