package service

import (
	"crypto/rand"
	"dst-admin-go/config/database"
	"dst-admin-go/config/global"
	"dst-admin-go/constant/dst"
	"dst-admin-go/model"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/utils/shellUtils"
	"dst-admin-go/vo"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type ClusterManager struct {
	DstHelper
	InitService
	ClusterService
}

func (c *ClusterManager) QueryCluster(ctx *gin.Context) {
	//获取查询参数
	clusterName := ctx.Query("clusterName")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	if page <= 0 {
		page = 1
	}
	if size < 0 {
		size = 10
	}
	db := database.DB

	if clusterName, isExist := ctx.GetQuery("clusterName"); isExist {
		db = db.Where("cluster_name LIKE ?", "%"+clusterName+"%")
	}

	db = db.Order("created_at desc").Limit(size).Offset((page - 1) * size)

	clusters := make([]model.Cluster, 0)

	if err := db.Find(&clusters).Error; err != nil {
		fmt.Println(err.Error())
	}

	var total int64
	db2 := database.DB
	if clusterName != "" {
		db2.Model(&model.Cluster{}).Where("clusterName like ?", "%"+clusterName+"%").Count(&total)
	} else {
		db2.Model(&model.Cluster{}).Count(&total)
	}
	totalPages := total / int64(size)
	if total%int64(size) != 0 {
		totalPages++
	}

	var clusterVOList = make([]vo.ClusterVO, len(clusters))
	var wg sync.WaitGroup
	wg.Add(len(clusters))
	for i, cluster := range clusters {
		go func(cluster model.Cluster, i int) {
			clusterVO := vo.ClusterVO{
				ClusterName:     cluster.ClusterName,
				Description:     cluster.Description,
				SteamCmd:        cluster.SteamCmd,
				ForceInstallDir: cluster.ForceInstallDir,
				Backup:          cluster.Backup,
				ModDownloadPath: cluster.ModDownloadPath,
				Uuid:            cluster.Uuid,
				Beta:            cluster.Beta,
				ID:              cluster.ID,
				CreatedAt:       cluster.CreatedAt,
				UpdatedAt:       cluster.UpdatedAt,
				//Master:          dst.Status(cluster.ClusterName, "Master"),
				//Caves:           dst.Status(cluster.ClusterName, "Caves"),
				Master: true,
				Caves:  true,
			}
			clusterIniPath := dst.GetClusterIniPath(clusterName)
			clusterIni := c.ReadClusterIniFile(clusterIniPath)
			name := clusterIni.ClusterName
			maxPlayers := clusterIni.MaxPlayers
			mode := clusterIni.GameMode
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
						info.Mode == mode &&
						int(info.Password) == hasPassword {
						clusterVO.RowId = info.Row
						clusterVO.Connected = int(info.Connected)
						clusterVO.MaxConnections = int(info.MaxConnect)
						clusterVO.Mode = info.Mode
						clusterVO.Mods = int(info.Mods)
						clusterVO.Season = info.Season
					}

				}
			}
			clusterVOList[i] = clusterVO
			wg.Done()
		}(cluster, i)
	}
	wg.Wait()
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: vo.Page{
			Data:       clusterVOList,
			Page:       page,
			Size:       size,
			Total:      total,
			TotalPages: totalPages,
		},
	})

}

func (c *ClusterManager) CreateCluster(cluster *model.Cluster) {

	if cluster.ClusterName == "" {
		log.Panicln("create cluster is error, cluster name is null")
	}
	if cluster.SteamCmd == "" {
		log.Panicln("create cluster is error, steamCmd is null")
	}
	if cluster.ForceInstallDir == "" {
		log.Panicln("create cluster is error, forceInstallDir is null")
	}
	db := database.DB
	cluster.Uuid = generateUUID()
	db.Create(&cluster)

	// 安装 dontstarve_dedicated_server
	log.Println("正在安装饥荒。。。。。。")
	if !fileUtils.Exists(cluster.ForceInstallDir) {
		// app_update 343050 beta updatebeta validate +quit
		cmd := "cd " + cluster.SteamCmd + " ; ./steamcmd.sh +login anonymous +force_install_dir " + cluster.ForceInstallDir + " +app_update 343050 validate +quit"
		output, err := shellUtils.Shell(cmd)
		if err != nil {
			log.Panicln("饥荒安装失败")
		}
		log.Println(output)
	}
	log.Println("饥荒安装完成！！！")
	// 创建世界
	c.InitCluster(cluster, global.CLUSTER_TOKEN)

}

func (c *ClusterManager) UpdateCluster(cluster *model.Cluster) {
	db := database.DB
	oldCluster := &model.Cluster{}
	db.Where("id = ?", cluster.ID).First(oldCluster)
	oldCluster.Description = cluster.Description
	//oldCluster.SteamCmd = cluster.SteamCmd
	//oldCluster.ForceInstallDir = cluster.ForceInstallDir
	db.Updates(oldCluster)
}

func (c *ClusterManager) DeleteCluster(id uint) (*model.Cluster, error) {
	db := database.DB

	cluster := model.Cluster{}
	result := db.Where("id = ?", id).Delete(&cluster)
	if result.Error != nil {
		return nil, result.Error
	}
	// TODO 删除集群 和 饥荒、备份、mod 下载

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
