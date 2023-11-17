package service

import (
	"crypto/rand"
	"dst-admin-go/config/database"
	"dst-admin-go/model"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/utils/dstUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sync"
)

type ClusterManager struct {
	InitService
	HomeService
	s GameService
	GameArchive
}

func (c *ClusterManager) QueryCluster(ctx *gin.Context) {
	//获取查询参数
	db := database.DB
	clusters := make([]model.Cluster, 0)
	if err := db.Find(&clusters).Error; err != nil {
		fmt.Println(err.Error())
	}

	var clusterVOList = make([]vo.ClusterVO, len(clusters))
	var wg sync.WaitGroup
	wg.Add(len(clusters))
	for i, cluster := range clusters {
		go func(cluster model.Cluster, i int) {
			clusterVO := vo.ClusterVO{
				Name:            cluster.Name,
				ClusterName:     cluster.ClusterName,
				Description:     cluster.Description,
				SteamCmd:        cluster.SteamCmd,
				ForceInstallDir: cluster.ForceInstallDir,
				Backup:          cluster.Backup,
				ModDownloadPath: cluster.ModDownloadPath,
				Beta:            cluster.Beta,
				Bin:             cluster.Bin,
				ID:              cluster.ID,
				CreatedAt:       cluster.CreatedAt,
				UpdatedAt:       cluster.UpdatedAt,
			}
			clusterIni := c.GetClusterIni(cluster.ClusterName)
			name := clusterIni.ClusterName
			maxPlayers := clusterIni.MaxPlayers
			// mode := clusterIni.GameMode
			password := clusterIni.ClusterPassword
			var hasPassword int
			if password == "" {
				hasPassword = 0
			} else {
				hasPassword = 1
			}
			// http 请求服务信息
			homeInfos := clusterUtils.GetDstServerInfo(name)
			if len(homeInfos) > 0 {
				for _, info := range homeInfos {
					if info.Name == name &&
						uint(info.MaxConnect) == maxPlayers &&
						int(info.Password) == hasPassword {
						clusterVO.RowId = info.Row
						clusterVO.Connected = int(info.Connected)
						clusterVO.MaxConnections = int(info.MaxConnect)
						clusterVO.Mode = info.Mode
						clusterVO.Mods = int(info.Mods)
						clusterVO.Season = info.Season
						clusterVO.Region = info.Region
					}

				}
			}
			clusterVO.GameArchive = c.GetGameArchive(clusterVO.ClusterName)
			clusterVO.Status = c.GetLevelStatus(clusterVO.ClusterName, "Master")
			clusterVOList[i] = clusterVO
			wg.Done()
		}(cluster, i)
	}
	wg.Wait()
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: clusterVOList,
	})

}

func (c *ClusterManager) CreateCluster(cluster *model.Cluster) {

	if cluster.Name == "" {
		log.Panicln("create cluster is error, name is null")
	}
	if cluster.ClusterName == "" {
		log.Panicln("create cluster is error, cluster name is null")
	}
	if cluster.SteamCmd == "" {
		log.Panicln("create cluster is error, steamCmd is null")
	}
	if cluster.ForceInstallDir == "" {
		log.Panicln("create cluster is error, forceInstallDir is null")
	}
	if cluster.ModDownloadPath == "" {
		cluster.ModDownloadPath = dstUtils.GetDoNotStarveTogetherPath()
	}

	db := database.DB
	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := db.Create(&cluster).Error
	if err != nil {
		if err.Error() == "Error 1062: Duplicate entry" {
			log.Panicln("唯一索引冲突！", err)
		}
		log.Panicln("创建房间失败！", err)
	}

	// 创建世界
	c.InitCluster(cluster, "")
	tx.Commit()

}

func (c *ClusterManager) UpdateCluster(cluster *model.Cluster) {
	db := database.DB
	oldCluster := &model.Cluster{}
	db.Where("ID = ?", cluster.ID).First(oldCluster)
	if oldCluster.ClusterName == "" {
		log.Panicln("未找到当前存档 clusterName: ", cluster.ClusterName, cluster.ID)
	}

	if cluster.SteamCmd != "" {
		oldCluster.SteamCmd = cluster.SteamCmd
	}
	if cluster.ForceInstallDir != "" {
		oldCluster.ForceInstallDir = cluster.ForceInstallDir
	}
	if cluster.Backup != "" {
		oldCluster.Backup = cluster.Backup
	}
	if cluster.ModDownloadPath != "" {
		oldCluster.ModDownloadPath = cluster.ModDownloadPath
	}
	oldCluster.Name = cluster.Name
	oldCluster.Bin = cluster.Bin
	db.Updates(oldCluster)
}

func (c *ClusterManager) DeleteCluster(clusterName string) (*model.Cluster, error) {
	// 停止服务
	c.s.StopGame(clusterName)
	db := database.DB
	cluster := model.Cluster{}
	result := db.Where("cluster_name = ?", clusterName).Unscoped().Delete(&cluster)

	if result.Error != nil {
		return nil, result.Error
	}

	// TODO 删除房间 和 饥荒、备份、mod 下载

	// 删除房间
	fileUtils.DeleteDir(dstUtils.GetClusterBasePath(clusterName))
	// 删除饥荒

	return &cluster, nil
}

func (c *ClusterManager) FindClusterByUuid(uuid string) *model.Cluster {
	db := database.DB
	cluster := &model.Cluster{}
	db.Where("uuid=?", uuid).First(cluster)
	return cluster
}

// 生成随机UUID
func generateUUID() string {
	// 生成随机字节序列
	var uuid [16]byte
	_, err := rand.Read(uuid[:])
	if err != nil {
		panic(err)
	}

	// 设置UUID版本和变体
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0xbf) | 0x80 // Variant 1

	// 将UUID转换为字符串并返回
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}
