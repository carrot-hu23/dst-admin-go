package api

import (
	"dst-admin-go/entity"
	"dst-admin-go/vo"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func NewProxy(ctx *gin.Context) {

	app := ctx.Param("name")

	route := entity.RoutingTable[app]

	url := ctx.Request.URL.Path
	url = strings.TrimPrefix(url, "/app/"+app)
	route.Proxy.Director = func(r *http.Request) {
		r.Header = ctx.Request.Header
		r.Host = route.Url.Host
		r.URL.Scheme = route.Url.Scheme
		r.URL.Host = route.Url.Host
		r.URL.Path = url
	}

	route.Proxy.ServeHTTP(ctx.Writer, ctx.Request)

}

func GetProxyEntity(ctx *gin.Context) {

	proxyEntities := make([]entity.Proxy, 0)

	db := entity.DB
	db.Find(&proxyEntities)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "Get proxy entity success",
		Data: proxyEntities,
	})

}

func CreateProxyEntity(ctx *gin.Context) {

	proxyEntity := entity.Proxy{}

	ctx.Bind(&proxyEntity)

	db := entity.DB
	db.Create(&proxyEntity)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "create proxy entity success",
	})

}

func UpdateProxyEntity(ctx *gin.Context) {

	proxyParam := vo.ProxyParam{}

	proxyEntity := entity.Proxy{}

	ctx.Bind(&proxyParam)
	fmt.Println(proxyParam)

	db := entity.DB

	db.Where("id=?", proxyParam.Id).First(&proxyEntity)

	proxyEntity.Name = proxyParam.Name
	proxyEntity.Description = proxyParam.Description
	proxyEntity.Ip = proxyParam.Ip
	proxyEntity.Port = proxyParam.Port

	db.Save(proxyEntity)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "update proxy entity success",
	})

}

func DeleteProxyEntity(ctx *gin.Context) {

	id := ctx.Query("id")

	proxyEntity := entity.Proxy{}
	db := entity.DB
	db.Where("id=?", id).Take(&proxyEntity)
	db.Delete(&proxyEntity)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "delete proxy entity success",
	})

}
