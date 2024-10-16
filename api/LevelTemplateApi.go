package api

import (
	"dst-admin-go/config/database"
	"dst-admin-go/model"
	"dst-admin-go/vo"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type LevelTemplateApi struct {
}

func (c *LevelTemplateApi) GetLevelTemplate(ctx *gin.Context) {

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	if page <= 0 {
		page = 1
	}
	if size < 0 {
		size = 10
	}
	db := database.DB
	db2 := database.DB
	if name, isExist := ctx.GetQuery("name"); isExist {
		db = db.Where("name LIKE ?", "%"+name+"%")
		db2 = db2.Where("name LIKE ?", "%"+name+"%")
	}
	db = db.Order("created_at desc").Limit(size).Offset((page - 1) * size)
	levelTemplates := make([]model.LevelTemplate, 0)

	if err := db.Find(&levelTemplates).Error; err != nil {
		fmt.Println(err.Error())
	}

	var total int64
	db2.Model(&model.LevelTemplate{}).Count(&total)

	totalPages := total / int64(size)
	if total%int64(size) != 0 {
		totalPages++
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: vo.Page{
			Data:       levelTemplates,
			Page:       page,
			Size:       size,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

func (c *LevelTemplateApi) CreateLevelTemplate(ctx *gin.Context) {

	template := model.LevelTemplate{}
	err := ctx.ShouldBind(&template)
	if err != nil {
		log.Panicln(err)
	}
	if template.Name == "" {
		log.Panicln("template name 不能为空")
	}
	if template.Modoverrides == "" {
		log.Panicln("Modoverrides 不能为空")
	}
	if template.Leveldataoverride1 == "" {
		log.Panicln("Leveldataoverride1 不能为空")
	}

	db := database.DB
	if template.ID == 0 {
		db.Create(&template)
	} else {
		db.Updates(&template)
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (c *LevelTemplateApi) UpdateTemplate(ctx *gin.Context) {

	template := model.LevelTemplate{}
	err := ctx.ShouldBind(&template)
	if err != nil {
		log.Panicln(err)
	}

	db := database.DB
	db.Updates(&template)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (c *LevelTemplateApi) DeleteTemplate(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Query("id"))
	db := database.DB
	template := model.LevelTemplate{}
	if err = db.First(&template, id).Error; err != nil {
		log.Panicln(err)
	}
	db.Delete(&template)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}
