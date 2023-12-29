package api

import (
	"dst-admin-go/autoCheck"
	"dst-admin-go/service"
	"dst-admin-go/utils/dstConfigUtils"
	"dst-admin-go/vo"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DstConfigApi struct {
}

var initEvnService = service.InitService{}

func (d *DstConfigApi) GetDstConfig(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: dstConfigUtils.GetDstConfig(),
	})
}

func (d *DstConfigApi) SaveDstConfig(ctx *gin.Context) {
	dstConfig := dstConfigUtils.NewDstConfig()
	err := ctx.Bind(dstConfig)
	if err != nil {
		log.Panicln(err)
	}
	dstConfigUtils.SaveDstConfig(dstConfig)
	initEvnService.InitBaseLevel(dstConfig, "默认初始化的世界", "pds-g^KU_qE7e8rv1^VVrVXd/01kBDicd7UO5LeL+uYZH1+geZlrutzItvOaw=", true)

	autoCheck.Manager.ReStart(dstConfig.Cluster)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "save dst_config success",
		Data: nil,
	})
}
