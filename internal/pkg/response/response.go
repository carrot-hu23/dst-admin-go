package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"` //提示代码
	Msg  string      `json:"msg"`  //提示信息
	Data interface{} `json:"data"` //数据
}

type Page struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	TotalPages int64       `json:"totalPages"`
	Page       int         `json:"page"`
	Size       int         `json:"size"`
}

// OkWithData 成功响应带数据
func OkWithData(data interface{}, ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "success",
		Data: data,
	})
}

// OkWithMessage 成功响应带消息
func OkWithMessage(message string, ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  message,
		Data: nil,
	})
}

// FailWithMessage 失败响应带消息
func FailWithMessage(message string, ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{
		Code: 500,
		Msg:  message,
		Data: nil,
	})
}

// OkWithPage 成功响应带分页数据
func OkWithPage(data interface{}, total, page, size int64, ctx *gin.Context) {
	totalPages := (total + size - 1) / size
	ctx.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "success",
		Data: Page{
			Data:       data,
			Total:      total,
			TotalPages: totalPages,
			Page:       int(page),
			Size:       int(size),
		},
	})
}
