package api

import (
	"dst-admin-go/config/database"
	"dst-admin-go/model"
	"dst-admin-go/vo"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type KvApi struct {
}

func (i *KvApi) GetKv(ctx *gin.Context) {

	key := ctx.Query("key")
	db := database.DB
	kv := model.KV{}
	db.Where("key = ?", key).Find(&kv)
	//if kv.Value != "Y" {

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: kv.Value,
	})
}

func (i *KvApi) SaveKv(ctx *gin.Context) {

	kv := model.KV{}
	err := ctx.ShouldBind(&kv)
	if err != nil {
		log.Panicln(err)
	}
	db := database.DB
	oldKv := model.KV{}
	db.Where("key = ?", kv.Key).Find(&oldKv)
	if oldKv.ID == 0 {
		db.Create(&kv)
	} else {
		oldKv.Value = kv.Value
		db.Save(&oldKv)
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: kv.Value,
	})
}
