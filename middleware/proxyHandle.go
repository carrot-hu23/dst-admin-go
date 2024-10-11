package middleware

import (
	"dst-admin-go/config/database"
	"dst-admin-go/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"strings"
)

// 白名单路径
var whiteList = map[string]bool{
	"/api/login":                           true,
	"/api/logout":                          true,
	"/ws":                                  true,
	"/api/bootstrap":                       true,
	"/api/init":                            true,
	"/api/install/steamcmd":                true,
	"/api/cluster":                         true,
	"/api/cluster/detail":                  true,
	"/api/cluster/container":               true,
	"/api/user/account":                    true,
	"/api/user/account/cluster":            true,
	"/api/user/account/cluster/permission": true,

	"/steam/dst/news":             true,
	"/api/game/system/info":       true,
	"/api/dst/home/server":        true,
	"/api/dst/home/server/detail": true,
}

func Proxy(c *gin.Context) {

	requestURL := c.Request.URL.String()
	urlParts := strings.Split(requestURL, "?")
	baseURL := urlParts[0]

	// 检查路径是否在白名单中
	if !strings.Contains(baseURL, "/api") || whiteList[baseURL] || strings.Contains(baseURL, "/api/cluster/detail") {
		// 如果在白名单中，直接处理请求
		c.Next()
		return
	}

	// 获取请求头中的cluster UUID
	clusterUUID := c.GetHeader("Cluster")

	// 根据UUID查询对应的服务器信息
	var cluster model.Cluster
	result := database.DB.Where("cluster_name = ?", clusterUUID).First(&cluster)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// 构建代理请求
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = "http" // 可根据需要修改为https
			req.URL.Host = fmt.Sprintf("%s:%d", cluster.Ip, cluster.Port)
			req.Host = req.URL.Host
			//req.Header.Set("Cookie", cache.GetToken(cluster))
		},
		ModifyResponse: func(resp *http.Response) error {
			if resp.StatusCode == http.StatusUnauthorized {
				//cache.TokenMemo.DeleteKey(cache.GenTokenKey(cluster))
				//resp.StatusCode = http.StatusGatewayTimeout
			}
			return nil
		},
	}

	// 执行代理请求
	proxy.ServeHTTP(c.Writer, c.Request)

	// 阻止 Gin 继续处理请求
	c.Abort()
}
