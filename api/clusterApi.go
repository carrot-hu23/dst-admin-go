package api

import (
	"dst-admin-go/config/global"
	"dst-admin-go/model"
	"dst-admin-go/service"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type ClusterApi struct{}

var clusterManager = service.ClusterManager{}

func (c *ClusterApi) GetClusterList(ctx *gin.Context) {
	clusterManager.QueryCluster(ctx, sessions)
}

func (c *ClusterApi) CreateCluster(ctx *gin.Context) {

	clusterModel := model.Cluster{}
	ctx.ShouldBind(&clusterModel)
	fmt.Printf("%v", clusterModel)

	if !fileUtils.Exists(clusterModel.SteamCmd) {
		log.Panicln("steamcmd 路径不存在 path: ", clusterModel.SteamCmd)
	}
	if !fileUtils.Exists(clusterModel.ForceInstallDir) {
		log.Panicln("饥荒 路径不存在 path: ", clusterModel.ForceInstallDir)
	}
	clusterManager.CreateCluster(&clusterModel)

	global.CollectMap.AddNewCollect(clusterModel.ClusterName)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})

}

func (c *ClusterApi) UpdateCluster(ctx *gin.Context) {
	clusterModel := model.Cluster{}
	err := ctx.ShouldBind(&clusterModel)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Printf("%v", clusterModel)

	if !fileUtils.Exists(clusterModel.SteamCmd) {
		log.Panicln("steamcmd 路径不存在 path: ", clusterModel.SteamCmd)
	}
	if !fileUtils.Exists(clusterModel.ForceInstallDir) {
		log.Panicln("饥荒 路径不存在 path: ", clusterModel.ForceInstallDir)
	}
	log.Println("clusterModel", clusterModel)
	clusterManager.UpdateCluster(&clusterModel)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})

}

func (c *ClusterApi) DeleteCluster(ctx *gin.Context) {

	clusterName := ctx.Query("clusterName")

	clusterModel, err := clusterManager.DeleteCluster(clusterName)
	log.Println("删除", clusterModel)
	if err != nil {
		log.Panicln("delete cluster error", err)
	}

	global.CollectMap.RemoveCollect(clusterModel.ClusterName)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})

}
