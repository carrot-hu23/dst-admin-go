package handler

import (
	"dst-admin-go/internal/model"
	"dst-admin-go/internal/pkg/response"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KvHandler struct {
	db *gorm.DB
}

func NewKvHandler(db *gorm.DB) *KvHandler {
	return &KvHandler{
		db: db,
	}
}

func (i *KvHandler) RegisterRoute(router *gin.RouterGroup) {
	router.GET("/api/kv", i.GetKv)
	router.POST("/api/kv", i.SaveKv)
}

// GetKv 生成 swagger 文档注释
// @Summary 获取kv值
// @Description 获取kv值
// @Tags kv
// @Param key query string true "key"
// @Success 200 {object} response.Response
// @Router /api/kv [get]
func (i *KvHandler) GetKv(ctx *gin.Context) {

	key := ctx.Query("key")
	db := i.db
	kv := model.KV{}
	db.Where("key = ?", key).Find(&kv)
	//if kv.Value != "Y" {

	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "success",
		Data: kv.Value,
	})
}

// SaveKv 生成 swagger 文档注释
// @Summary 保存kv值
// @Description 保存kv值
// @Tags kv
// @Param key formData string true "key"
// @Param value formData string true "value"
// @Success 200 {object} response.Response
// @Router /api/kv [post]
func (i *KvHandler) SaveKv(ctx *gin.Context) {

	kv := model.KV{}
	err := ctx.ShouldBind(&kv)
	if err != nil {
		log.Panicln(err)
	}
	db := i.db
	oldKv := model.KV{}
	db.Where("key = ?", kv.Key).Find(&oldKv)
	if oldKv.ID == 0 {
		db.Create(&kv)
	} else {
		oldKv.Value = kv.Value
		db.Save(&oldKv)
	}

	ctx.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "success",
		Data: kv.Value,
	})
}
