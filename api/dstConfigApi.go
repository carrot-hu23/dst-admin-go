package api

import (
	"dst-admin-go/service"
	"dst-admin-go/utils/dstConfigUtils"
	"dst-admin-go/vo"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDstConfig(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: dstConfigUtils.GetDstConfig(),
	})
}

func SaveDstConfig(ctx *gin.Context) {
	dstConfig := dstConfigUtils.NewDstConfig()
	ctx.Bind(dstConfig)
	dstConfigUtils.SaveDstConfig(dstConfig)
	service.InitBaseLevel(dstConfig, "test", "pds-g^KU_qE7e8rv1^VVrVXd/01kBDicd7UO5LeL+uYZH1+geZlrutzItvOaw=", true)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "save dst_config success",
		Data: nil,
	})
}
