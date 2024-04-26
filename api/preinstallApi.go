package api

import (
	"dst-admin-go/autoCheck"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/utils/dstUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type PreinstallApi struct{}

func (p *PreinstallApi) UsePreinstall(ctx *gin.Context) {
	lock.Lock()
	defer lock.Unlock()
	cluster := clusterUtils.GetClusterFromGin(ctx)

	// TODO 关闭之前的世界
	gameService.StopGame(cluster.ClusterName)

	// 创建备份
	backupService.CreateBackup(cluster.ClusterName, "")

	name := ctx.DefaultQuery("name", "default")
	log.Println(name)

	if !fileUtils.Exists("./static/preinstall/" + name) {
		log.Panicln("./static/preinstall/"+name, "不存在，请先添加")
	}

	path := dstUtils.GetClusterBasePath(cluster.ClusterName)
	bakpath := filepath.Join(dstUtils.GetKleiDstPath(), "bak")
	//删除上次的备份目录 防止更名失败
	err := fileUtils.DeleteDir(bakpath)
	if err != nil {
		log.Panicln(err)
	}

	//先重命名但是不删除
	fileUtils.Rename(path, bakpath)

	err = fileUtils.Copy("./static/preinstall/"+name, filepath.Join(dstUtils.GetKleiDstPath()))
	if err != nil {
		log.Panicln(err)
	}
	newpath := dstUtils.GetClusterBasePath(cluster.ClusterName)
	fileUtils.Rename(filepath.Join(dstUtils.GetKleiDstPath(), name), newpath)
	//还原cluster_token 等
	tocopy := []string{"adminlist.txt", "blocklist.txt", "cluster_token.txt", "whitelist.txt"}
	for _, s := range tocopy {
		if fileUtils.Exists(filepath.Join(bakpath, s)) {
			err = fileUtils.Copy(filepath.Join(bakpath, s), newpath)
			if err != nil {
				log.Panicln(err)
			}
		}
	}
	err = fileUtils.DeleteDir(bakpath)
	if err != nil {
		log.Panicln(err)
	}
	// TODO 宕机恢复重新读取
	autoCheck.Manager.ReStart(cluster.ClusterName)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})

}
