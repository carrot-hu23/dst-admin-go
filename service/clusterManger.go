package service

import (
	"crypto/rand"
	"dst-admin-go/config/database"
	"dst-admin-go/config/global"
	"dst-admin-go/model"
	"dst-admin-go/utils/dstUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

type ClusterManager struct {
	InitService
	HomeService
	s GameService
	GameArchive
	RemoteService
}

func (c *ClusterManager) QueryCluster(ctx *gin.Context) {
	//获取查询参数
	db := database.DB
	clusters := make([]model.Cluster, 0)
	session := sessions.Default(ctx)
	role := session.Get("role")
	userId := session.Get("userId")
	log.Println("role", role, "userId", userId)
	if role == "admin" {
		if err := db.Find(&clusters).Error; err != nil {
			fmt.Println(err.Error())
		}
	} else {
		db2 := database.DB
		userClusterList := []model.UserCluster{}
		db2.Where("user_id = ?", userId).Find(&userClusterList)
		ids := []int{}
		for i := range userClusterList {
			ids = append(ids, userClusterList[i].ClusterId)
		}
		db.Where("id in ?", ids).Find(&clusters)
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
				Ugc_directory:   cluster.Ugc_directory,
				LevelType:       cluster.LevelType,
				ClusterType:     cluster.ClusterType,
				Ip:              cluster.Ip,
				Port:            cluster.Port,
				Username:        cluster.Username,
				ClusterPassword: cluster.Password,
			}
			if cluster.ClusterType == "远程" {
				// TODO 增加xinxi
				clusterVO.GameArchive = c.GetRemoteGameArchive(cluster)
				clusterVO.Status = c.GetRemoteLevelStatus(cluster)
			} else {
				clusterVO.GameArchive = c.GetGameArchive(clusterVO.ClusterName)
				clusterVO.Status = c.GetLevelStatus(clusterVO.ClusterName, "Master")
			}
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
	if cluster.ClusterType == "本地" {
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
			p := filepath.Join(dstUtils.GetDoNotStarveTogetherPath(), "mod_download")
			fileUtils.CreateDirIfNotExists(p)
			cluster.ModDownloadPath = p
		}
	}

	db := database.DB
	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	if len(cluster.RemoteClusterNameList) > 0 {
		if len(cluster.RemoteClusterNameList) > 0 {
			var clusterList []model.Cluster
			remoteClusterNameList := cluster.RemoteClusterNameList
			for i := range remoteClusterNameList {
				remoteClusterName := remoteClusterNameList[i]
				newUUID, _ := uuid.NewUUID()
				newCluster := model.Cluster{
					ClusterName:       newUUID.String(),
					RemoteClusterName: remoteClusterName,
					Ip:                cluster.Ip,
					Port:              cluster.Port,
					Username:          cluster.Username,
					Password:          cluster.Password,
					ClusterType:       cluster.ClusterType,
					Name:              cluster.Name,
				}
				clusterList = append(clusterList, newCluster)
			}
			tx.Create(&clusterList)
		}
		return
	}

	err := tx.Create(&cluster).Error
	if err != nil {
		if err.Error() == "Error 1062: Duplicate entry" {
			log.Panicln("唯一索引冲突！", err)
		}
		log.Panicln("创建房间失败！", err)
	}
	if cluster.ClusterType != "远程" {
		// 创建世界
		c.InitCluster2(cluster, global.Config.Token)
	}
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
	oldCluster.Ugc_directory = cluster.Ugc_directory
	db.Updates(oldCluster)

	if cluster.Ugc_directory == "" {
		db.Model(&model.Cluster{}).Where("ID = ?", cluster.ID).UpdateColumn("ugc_directory", "")
	}

}

func (c *ClusterManager) DeleteCluster(clusterName string) (*model.Cluster, error) {

	if clusterName == "" {
		log.Panicln("cluster is not allow null")
	}

	// 停止服务
	c.s.StopGame(clusterName)
	db := database.DB
	cluster := model.Cluster{}
	result := db.Where("cluster_name = ?", clusterName).Unscoped().Delete(&cluster)

	if result.Error != nil {
		return nil, result.Error
	}

	// TODO 删除房间 和 饥荒、备份、mod 下载
	if cluster.ClusterType != "远程" {
		// 删除房间
		fileUtils.DeleteDir(dstUtils.GetClusterBasePath(clusterName))
		// 删除饥荒
	}

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
