package api

import (
	"dst-admin-go/config/database"
	"dst-admin-go/config/dockerClient"
	"dst-admin-go/model"
	"dst-admin-go/service"
	"dst-admin-go/session"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo"
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ClusterApi struct{}

var clusterManager = service.ClusterManager{}
var portInfoService = service.PortInfoService{}

func (c *ClusterApi) GetClusterList(ctx *gin.Context) {
	clusterManager.QueryCluster(ctx, sessions)
}

func checkAdmin(ctx *gin.Context, sessions *session.Manager) {
	s := sessions.Start(ctx.Writer, ctx.Request)
	role := s.Get("role")
	if role != "admin" {
		log.Panicln("你无权限操作")
	}
}

func (c *ClusterApi) CreateCluster(ctx *gin.Context) {

	checkAdmin(ctx, sessions)

	clusterModel := model.Cluster{}
	err := ctx.ShouldBind(&clusterModel)
	if err != nil {
		log.Panicln(err)
	}
	if clusterModel.Day == 0 {
		log.Panicln("过期时间不能为0")
	}
	if clusterModel.LevelNum == 0 {
		log.Panicln("世界层数不能为0")
	}
	if clusterModel.MaxPlayers == 0 {
		log.Panicln("玩家人数不能为0")
	}
	if clusterModel.Core == 0 {
		log.Panicln("cpu核数不能为0")
	}
	if clusterModel.Memory == 0 {
		log.Panicln("内存不能为0")
	}
	if clusterModel.MaxBackup == 0 {
		log.Panicln("最大备份数量不能为0")
	}
	fmt.Printf("%v", clusterModel)

	var clusterList []model.Cluster

	zone, ok := dockerClient.Zone(clusterModel.ZoneCode)
	if !ok {
		log.Panicln("未找到当前 ", clusterModel.ZoneCode)
	}

	// 批量创建
	quantity := clusterModel.Quantity
	for i := 0; i < quantity; i++ {
		cluster := model.Cluster{
			LevelNum:   clusterModel.LevelNum,
			MaxPlayers: clusterModel.MaxPlayers,
			MaxBackup:  clusterModel.MaxBackup,
			Memory:     clusterModel.Memory,
			Core:       clusterModel.Core,
			Disk:       clusterModel.Disk,
			Day:        clusterModel.Day,
			Name:       fmt.Sprintf("%s-%d", clusterModel.Name, i+1),
			Image:      clusterModel.Image,
			ZoneCode:   zone.ZoneCode,
			ZoneName:   zone.Name,
			Ip:         zone.Ip,
		}
		log.Println("正在创建cluster", cluster)
		e := clusterManager.CreateCluster(&cluster)
		if e != nil {
			log.Panicln(e)
		}
		clusterList = append(clusterList, cluster)
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: clusterList,
	})

}

func (c *ClusterApi) UpdateCluster(ctx *gin.Context) {
	clusterModel := model.Cluster{}
	err := ctx.ShouldBind(&clusterModel)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Printf("%v", clusterModel)
	log.Println("clusterModel", clusterModel)
	clusterManager.UpdateCluster(&clusterModel)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})

}

func (c *ClusterApi) DeleteCluster(ctx *gin.Context) {

	checkAdmin(ctx, sessions)

	clusterName := ctx.Query("clusterName")

	clusterModel, err := clusterManager.DeleteCluster(clusterName)
	log.Println("删除", clusterModel)
	if err != nil {
		log.Panicln("delete cluster error", err)
	}
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})

}

func (c *ClusterApi) RestartCluster(ctx *gin.Context) {

	clusterName := ctx.Query("clusterName")

	err := clusterManager.RestartContainer(clusterName)
	log.Println("重启", clusterName)
	if err != nil {
		log.Panicln("restart cluster error", err)
	}
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})

}

func (c *ClusterApi) GetCluster(ctx *gin.Context) {

	clusterName := ctx.Param("id")
	fmt.Printf("%s", clusterName)

	db := database.DB
	var cluster model.Cluster
	db.Where("cluster_name = ?", clusterName).Find(&cluster)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: cluster,
	})
}

func (c *ClusterApi) UpdateClusterContainer(ctx *gin.Context) {

	checkAdmin(ctx, sessions)

	var payload struct {
		ClusterName string `json:"ClusterName"`
		Day         int64  `json:"day"`
	}
	err := ctx.ShouldBind(&payload)
	if err != nil {
		log.Panicln(err)
	}

	db := database.DB
	var cluster model.Cluster
	db.Where("cluster_name = ?", payload.ClusterName).Find(&cluster)

	cluster.Day = cluster.Day + payload.Day
	cluster.ExpireTime = cluster.ExpireTime + payload.Day*24*60*60

	db.Save(&cluster)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: cluster,
	})

}

func (c *ClusterApi) BindCluster(ctx *gin.Context) {
	var payload struct {
		ClusterName string `json:"ClusterName"`
		Username    string `json:"username"`
		DisplayName string `json:"displayName"`
		Password    string `json:"password"`
		Description string `json:"description"`
		PhotoURL    string `json:"photoURL"`
	}
	err := ctx.ShouldBind(&payload)
	if err != nil {
		log.Panicln(err)
	}
	log.Println("激活卡密", payload)

	db1 := database.DB
	oldUser := model.User{}
	db1.Where("username=?", payload.Username).First(&oldUser)
	if oldUser.ID != 0 {
		ctx.JSON(http.StatusOK, vo.Response{
			Code: 531,
			Msg:  "用户名重复了,请换一个",
			Data: nil,
		})
		return
	}

	db := database.DB
	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			ctx.JSON(http.StatusOK, vo.Response{
				Code: 500,
				Msg:  "绑定失败",
				Data: nil,
			})
		}
	}()

	// 绑定
	cluster := model.Cluster{}
	tx.Where("cluster_name = ?", payload.ClusterName).Find(&cluster)

	if cluster.Activate {
		log.Panicln("当前卡密已激活，无法绑定")
	}

	// 创建用户
	user := model.User{
		Username:    payload.Username,
		Password:    payload.Password,
		DisplayName: payload.DisplayName,
		PhotoURL:    payload.PhotoURL,
	}

	tx.Create(&user)
	log.Println("创建用户成功", user)

	userCluster := model.UserCluster{}
	userCluster.ClusterId = int(cluster.ID)
	userCluster.UserId = int(user.ID)

	log.Println("正在绑定", userCluster)
	tx.Create(&userCluster)

	portCount := 1 + 1 + cluster.LevelNum + 5
	ports, err := portInfoService.GetAvailablePorts(cluster.ZoneCode, portCount)
	if err != nil {
		log.Panicln("获取端口失败")
	}
	cluster.Port = ports[0]
	cluster.MasterPort = ports[1]

	// 激活卡密
	containerId, err := clusterManager.CreateContainer(cluster)
	if err != nil {
		log.Panicln(err)
	}
	cluster.Activate = true
	cluster.ContainerId = containerId
	cluster.Expired = false
	cluster.ExpireTime = time.Now().Add(time.Duration(cluster.Day) * time.Hour * 24).Unix()
	tx.Save(&cluster)

	err = portInfoService.SaveAvailablePort(tx, cluster.ZoneCode, containerId, ports)
	if err != nil {
		log.Panicln(err)
	}
	tx.Commit()

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})

}

func (c *ClusterApi) GetKamiList(ctx *gin.Context) {

	page, _ := strconv.Atoi(ctx.DefaultQuery("current", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	if page <= 0 {
		page = 1
	}
	if size < 0 {
		size = 10
	}

	db := database.DB
	db2 := database.DB

	if zoneCode, isExist := ctx.GetQuery("zoneCode"); isExist {
		db = db.Where("zone_code = ?", zoneCode)
		db2 = db2.Where("zone_code = ?", zoneCode)
	}

	if levelNum, isExist := ctx.GetQuery("levelNum"); isExist {
		intValue, _ := strconv.Atoi(levelNum)
		db = db.Where("level_num = ?", intValue)
		db2 = db2.Where("level_num = ?", intValue)
	}
	if maxPlayers, isExist := ctx.GetQuery("maxPlayers"); isExist {
		intValue, _ := strconv.Atoi(maxPlayers)
		db = db.Where("max_players = ?", intValue)
		db2 = db2.Where("max_players = ?", intValue)
	}
	if core, isExist := ctx.GetQuery("core"); isExist {
		intValue, _ := strconv.Atoi(core)
		db = db.Where("core = ?", intValue)
		db2 = db2.Where("core = ?", intValue)
	}
	if memory, isExist := ctx.GetQuery("memory"); isExist {
		intValue, _ := strconv.Atoi(memory)
		db = db.Where("memory = ?", intValue)
		db2 = db2.Where("memory = ?", intValue)
	}

	db = db.Where("activate = ?", false)
	db2 = db2.Where("activate = ?", false)

	db = db.Order("created_at desc").Limit(size).Offset((page - 1) * size)
	clusters := make([]model.Cluster, 0)
	db.Find(&clusters)

	var total int64
	db2.Model(&model.Cluster{}).Count(&total)
	totalPages := total / int64(size)
	if total%int64(size) != 0 {
		totalPages++
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: vo.Page{
			Data:       clusters,
			Page:       page,
			Size:       size,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

func (c *ClusterApi) ExportKamiList(ctx *gin.Context) {

	// 获取当前时间并格式化为字符串
	currentTime := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("%s-dst-kami.csv", currentTime)
	// 设置响应头信息，指定导出的文件名
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment;filename=%s", filename))
	ctx.Header("Content-Type", "text/csv; charset=UTF-8")
	ctx.Header("Content-Transfer-Encoding", "binary")

	// 在响应体中写入 UTF-8 BOM 头
	ctx.Writer.Write([]byte{0xEF, 0xBB, 0xBF})

	// 创建 CSV Writer，将数据写入 Response
	writer := csv.NewWriter(ctx.Writer)
	defer writer.Flush()

	db := database.DB
	if zoneCode, isExist := ctx.GetQuery("zoneCode"); isExist {
		db = db.Where("zone_code = ?", zoneCode)
	}

	if levelNum, isExist := ctx.GetQuery("levelNum"); isExist {
		intValue, _ := strconv.Atoi(levelNum)
		db = db.Where("level_num = ?", intValue)
	}
	if maxPlayers, isExist := ctx.GetQuery("maxPlayers"); isExist {
		intValue, _ := strconv.Atoi(maxPlayers)
		db = db.Where("max_players = ?", intValue)
	}
	if core, isExist := ctx.GetQuery("core"); isExist {
		intValue, _ := strconv.Atoi(core)
		db = db.Where("core = ?", intValue)
	}
	if memory, isExist := ctx.GetQuery("memory"); isExist {
		intValue, _ := strconv.Atoi(memory)
		db = db.Where("memory = ?", intValue)
	}
	db = db.Where("activate = ?", false)
	db = db.Order("created_at desc")
	clusters := make([]model.Cluster, 0)
	db.Find(&clusters)

	// 写入 CSV 数据
	data := [][]string{
		{"卡密", "区域", "内存(GB)", "核数", "世界层数", "最大玩家", "天数"},
	}

	for i := range clusters {
		data = append(data, []string{
			clusters[i].Uuid,
			clusters[i].ZoneName,
			fmt.Sprintf("%d", clusters[i].Memory),
			fmt.Sprintf("%d", clusters[i].Core),
			fmt.Sprintf("%d", clusters[i].LevelNum),
			fmt.Sprintf("%d", clusters[i].MaxPlayers),
			fmt.Sprintf("%d", clusters[i].Day),
		})
	}

	// 遍历数据并写入 CSV
	for _, row := range data {
		if err := writer.Write(row); err != nil {
			ctx.String(http.StatusInternalServerError, "Error writing CSV: %v", err)
			return
		}
	}

	// 添加 Boom 头
	ctx.Header("Content-Security-Policy", "default-src 'none'; style-src 'self'; font-src 'self';")
}

func getStartPort() int64 {
	version, err := fileUtils.ReadFile("./startPort")
	if err != nil {
		log.Println(err)
		return 20000
	}
	version = strings.Replace(version, "\n", "", -1)
	l, err := strconv.ParseInt(version, 10, 64)
	if err != nil {
		log.Println(err)
		return 20000
	}
	return l
}

func saveEndPort(portEnd int64) {
	fileUtils.WriterTXT("./startPort", strconv.Itoa(int(portEnd)))
}
