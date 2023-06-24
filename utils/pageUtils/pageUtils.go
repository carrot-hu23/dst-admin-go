package pageUtils

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func RequestPage(ctx *gin.Context) (int, int) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	if page <= 0 {
		page = 1
	}
	if size < 0 {
		size = 10
	}
	return page, size
}
