package context

import (
	"dst-admin-go/internal/service/dstConfig"

	"github.com/gin-gonic/gin"
)

const (
	clusterNameKey = "cluster_name"
	dstConfigKey   = "dst_config"
)

// GetClusterName 从 gin.Context 获取集群名称
// 需要配合 ClusterMiddleware 使用
func GetClusterName(c *gin.Context) string {
	if value, exists := c.Get(clusterNameKey); exists {
		if clusterName, ok := value.(string); ok {
			return clusterName
		}
	}
	return ""
}

// GetDstConfig 从 gin.Context 获取 DstConfig 对象
// 需要配合 ClusterMiddleware 使用
func GetDstConfig(c *gin.Context) *dstConfig.DstConfig {
	if value, exists := c.Get(dstConfigKey); exists {
		if config, ok := value.(dstConfig.DstConfig); ok {
			return &config
		}
	}
	return nil
}
