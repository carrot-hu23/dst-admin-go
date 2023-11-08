package api

import (
	"dst-admin-go/autoCheck"
	"dst-admin-go/constant"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/utils/dstUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"path/filepath"
)

type PreinstallApi struct{}

func (p *PreinstallApi) UsePreinstall(ctx *gin.Context) {
	lock.Lock()
	defer lock.Unlock()
	cluster := clusterUtils.GetClusterFromGin(ctx)

	// TODO 关闭之前的世界
	gameService.StopGame(cluster.ClusterName)

	// 创建备份
	backupService.CreateBackup(ctx, "")

	name := ctx.DefaultQuery("name", "default")
	log.Println(name)

	if !fileUtils.Exists("./static/preinstall/" + name) {
		log.Panicln("./static/preinstall/"+name, "不存在，请先添加")
	}
	err := fileUtils.DeleteDir(dstUtils.GetClusterBasePath(cluster.ClusterName))
	if err != nil {
		log.Panicln(err)
	}
	err = fileUtils.Copy("./static/preinstall/"+name, filepath.Join(constant.HOME_PATH, ".klei", "DoNotStarveTogether"))
	if err != nil {
		log.Panicln(err)
	}
	fileUtils.Rename(filepath.Join(constant.HOME_PATH, ".klei", "DoNotStarveTogether", name), dstUtils.GetClusterBasePath(cluster.ClusterName))

	// TODO 宕机恢复重新读取
	autoCheck.Manager.ReStart()

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})

}
