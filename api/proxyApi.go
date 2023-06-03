package api

import (
	"dst-admin-go/config/database"
	"dst-admin-go/config/global"
	"dst-admin-go/model"
	"dst-admin-go/vo"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ProxyApi struct {
}

func (p *ProxyApi) NewProxy(ctx *gin.Context) {

	app := ctx.Param("name")

	route := global.RoutingTable[app]

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

func (p *ProxyApi) GetProxyEntity(ctx *gin.Context) {

	proxyEntities := make([]model.Proxy, 0)

	db := database.DB
	db.Find(&proxyEntities)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "Get proxy model success",
		Data: proxyEntities,
	})

}

func (p *ProxyApi) CreateProxyEntity(ctx *gin.Context) {

	proxyEntity := model.Proxy{}

	ctx.Bind(&proxyEntity)

	db := database.DB
	db.Create(&proxyEntity)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "create proxy model success",
	})

}

func (p *ProxyApi) UpdateProxyEntity(ctx *gin.Context) {

	proxyParam := vo.ProxyParam{}

	proxyEntity := model.Proxy{}

	ctx.Bind(&proxyParam)
	fmt.Println(proxyParam)

	db := database.DB

	db.Where("id=?", proxyParam.Id).First(&proxyEntity)

	proxyEntity.Name = proxyParam.Name
	proxyEntity.Description = proxyParam.Description
	proxyEntity.Ip = proxyParam.Ip
	proxyEntity.Port = proxyParam.Port

	db.Save(proxyEntity)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "update proxy model success",
	})

}

func (p *ProxyApi) DeleteProxyEntity(ctx *gin.Context) {

	id := ctx.Query("id")

	proxyEntity := model.Proxy{}
	db := database.DB
	db.Where("id=?", id).Take(&proxyEntity)
	db.Delete(&proxyEntity)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "delete proxy model success",
	})

}
