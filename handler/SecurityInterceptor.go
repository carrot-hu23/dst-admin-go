package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ShellInjectionInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求参数
		params := c.Request.URL.Query()
		for _, param := range params {
			// 检查参数是否包含 Shell 特殊字符
			if strings.ContainsAny(param[0], "|&;<>()$`\\\"'*?#[]{}~=") {
				// 如果包含，则返回错误信息
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid characters in parameter"})
				c.Abort()
				return
			}
			// 检查参数是否包含 SQL 注入特殊字符
			if strings.ContainsAny(param[0], "'\";\\") {
				// 如果包含，则返回错误信息
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid characters in parameter"})
				c.Abort()
				return
			}
			// 检查参数是否包含系统文件路径或Linux shadow文件路径
			if strings.Contains(param[0], "/etc/") || strings.Contains(param[0], "/usr/bin/") || strings.Contains(param[0], "/etc/shadow") || strings.Contains(param[0], "../") || strings.Contains(param[0], "..") {
				// 如果包含，则返回错误信息
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid characters in parameter"})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
