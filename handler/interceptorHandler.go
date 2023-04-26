package handler

import (
	"dst-admin-go/api"
	"dst-admin-go/constant"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// 拦截检查是否安装 dst steam cmd
func CheckDstHandler() gin.HandlerFunc {

	return func(c *gin.Context) {

		//request := c.Request
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
			// return可省略, 只要前面执行Abort()就可以让后面的handler函数不再执行
			return
		}
	}

}

func checkIsInstallDst(path string) bool {
	// for _, value := range blacklist {
	// 	if value == path {
	// 		return true
	// 	}
	// }
	if filter(whilelist, path) {
		return true
	}

	dst_path := constant.DST_INSTALL_DIR + "/bin" + constant.SINGLE_SLASH + constant.DST_START_PROGRAM
	log.Println("dst_path", dst_path)

	return fileUtils.Exists(dst_path)
}

var whilelist = []string{"/api/login", "/api/logout", "/ws", "/api/init"}

func filter(s []string, str string) bool {
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
func Authentucation() gin.HandlerFunc {

	return func(c *gin.Context) {

		path := c.Request.URL.Path
		if filter(whilelist, path) {
			c.Next()
			return
		} else {
			session := api.Sessions().Start(c.Writer, c.Request)
			//fmt.Sprintf("%v", session.Get("username"))
			cookieName := session.Get("username")
			sessionID := url.QueryEscape(session.SessionID())
			log.Println("cookiName: " + fmt.Sprintf("%v", session.Get("username")))
			log.Println("sessionID: " + sessionID)
			if cookieName == nil {
				// c.Abort()
				// c.JSON(http.StatusBadGateway, vo.Response{
				// 	Code: 401,
				// 	Msg:  "Please login",
				// })
				// 如果用户未登录，返回 HTTP 401
				c.AbortWithStatus(http.StatusUnauthorized)
			} else {
				c.Next()
			}
		}

	}

}
