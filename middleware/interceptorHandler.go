package middleware

import (
	"dst-admin-go/api"
	"dst-admin-go/vo"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	whitelist = []string{"/api/login", "/api/logout", "/ws", "/api/bootstrap", "/api/init", "/api/install/steamcmd"}
)

// 拦截检查是否安装 dst steam cmd
func CheckDstHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if checkIsInstallDst(path) {
			// 验证通过，会继续访问下一个中间件
			c.Next()
		} else {
			// 验证不通过，不再调用后续的函数处理
			c.Abort()
			c.JSON(http.StatusBadGateway, vo.Response{
				Code: 501,
				Msg:  "Sorry, you haven't installed DST Stream CMD",
			})
			return
		}
	}
}

func checkIsInstallDst(path string) bool {

	if apiFilter(whitelist, path) {
		return true
	}

	//dstPath := constant.DST_INSTALL_DIR + "/bin" + constant.SINGLE_SLASH + constant.DST_START_PROGRAM
	//log.Println("dst_path", dstPath)
	//
	//return fileUtils.Exists(dstPath)
	return true
}

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
			session := api.Sessions().Start(c.Writer, c.Request)
			cookieName := session.Get("username")
			// sessionID := url.QueryEscape(session.SessionID())
			if cookieName == nil {
				// 如果用户未登录，返回 HTTP 401
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
