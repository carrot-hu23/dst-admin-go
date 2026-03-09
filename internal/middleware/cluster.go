package middleware

import (
	"dst-admin-go/internal/pkg/response"
	"dst-admin-go/internal/service/dstConfig"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	clusterNameKey = "cluster_name"
	dstConfigKey   = "dst_config"
)

// ClusterMiddleware 从 HTTP Header 解析 cluster 名称并加载配置到 context
func ClusterMiddleware(dstConfigService dstConfig.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		path := c.Request.URL.Path
		if apiFilter(whitelist, path) {
			c.Next()
			return
		}

		// 从 Header 获取集群名称
		clusterName := c.GetHeader("Cluster")

		// 查询集群配置
		config, err := dstConfigService.GetDstConfig(clusterName)
		if err != nil {
			c.JSON(http.StatusNotFound, response.Response{
				Code: 500,
				Msg:  "集群配置不存在: " + clusterName,
			})
			c.Abort()
			return
		}

		// 将集群名称和配置注入到 context
		c.Set(clusterNameKey, config.Cluster)
		c.Set(dstConfigKey, config)

		c.Next()
	}
}
