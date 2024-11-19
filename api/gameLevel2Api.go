package api

import (
	"dst-admin-go/constant/consts"
	"dst-admin-go/model"
	"dst-admin-go/service"
	"dst-admin-go/service/autoCheck"
	"dst-admin-go/utils/clusterUtils"
	"dst-admin-go/utils/dstConfigUtils"
	"dst-admin-go/utils/shellUtils"
	"dst-admin-go/vo"
	"dst-admin-go/vo/level"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var gameLevel2Service = service.GameLevel2Service{}
var homeService = service.HomeService{}

type GameLevel2Api struct{}

type ClusterIni struct {
	Cluster *level.ClusterIni `json:"cluster"`
	Token   string            `json:"token"`
}

type LevelInfo struct {
	Ps                *vo.DstPsVo      `json:"Ps"`
	Status            bool             `json:"status"`
	LevelName         string           `json:"levelName"`
	IsMaster          bool             `json:"is_master"`
	Uuid              string           `json:"uuid"`
	Leveldataoverride string           `json:"leveldataoverride"`
	Modoverrides      string           `json:"modoverrides"`
	ServerIni         *level.ServerIni `json:"server_ini"`
}

// GetStatus 获取世界状态
func (g *GameLevel2Api) GetStatus(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	levelList := gameLevel2Service.GetLevelList(clusterName)
	length := len(levelList)
	result := make([]LevelInfo, length)

	if runtime.GOOS == "windows" {

		var wg sync.WaitGroup
		wg.Add(length)
		for i := range levelList {
			go func(index int) {
				defer func() {
					wg.Done()
					if r := recover(); r != nil {

					}
				}()
				world := levelList[index]
				ps := gameService.PsAuxSpecified(clusterName, world.Uuid)
				status := gameService.GetLevelStatus(clusterName, world.Uuid)
				result[index] = LevelInfo{
					Ps:                ps,
					Status:            status,
					LevelName:         world.LevelName,
					IsMaster:          world.IsMaster,
					Uuid:              world.Uuid,
					Leveldataoverride: world.Leveldataoverride,
					Modoverrides:      world.Modoverrides,
					ServerIni:         world.ServerIni,
				}
			}(i)
		}
		wg.Wait()
		ctx.JSON(http.StatusOK, vo.Response{
			Code: 200,
			Msg:  "success",
			Data: result,
		})
	} else {
		for i := range levelList {
			world := levelList[i]
			result[i] = LevelInfo{
				Ps:                &vo.DstPsVo{},
				Status:            false,
				LevelName:         world.LevelName,
				IsMaster:          world.IsMaster,
				Uuid:              world.Uuid,
				Leveldataoverride: world.Leveldataoverride,
				Modoverrides:      world.Modoverrides,
				ServerIni:         world.ServerIni,
			}
		}

		cmd := "ps -aux | grep -v grep | grep -v tail | grep -v SCREEN | grep " + clusterName + " |awk '{print $3, $4, $5, $6,$16}'"
		info, err := shellUtils.Shell(cmd)
		if err != nil {
			log.Println(cmd + " error: " + err.Error())
		} else {
			lines := strings.Split(info, "\n")
			for lineIndex := range lines {
				dstPsVo := vo.NewDstPsVo()
				arr := strings.Split(lines[lineIndex], " ")
				//for i := range arr {
				//	log.Println(i, arr[i])
				//}
				if len(arr) > 4 {
					dstPsVo.CpuUage = strings.Replace(arr[0], "\n", "", -1)
					dstPsVo.MemUage = strings.Replace(arr[1], "\n", "", -1)
					dstPsVo.VSZ = strings.Replace(arr[2], "\n", "", -1)
					dstPsVo.RSS = strings.Replace(arr[3], "\n", "", -1)
					for i := range result {
						levelName := result[i].Uuid
						if strings.Contains(arr[4], levelName) {
							result[i].Ps = dstPsVo
							result[i].Status = true
						}
					}
				}

			}
		}

		ctx.JSON(http.StatusOK, vo.Response{
			Code: 200,
			Msg:  "success",
			Data: result,
		})
	}

}

func startBefore(ctx *gin.Context) {
	copyOsFile()
}

func copyOsFile() {
	config := dstConfigUtils.GetDstConfig()
	// 定义路径
	srcFile := filepath.Join(config.Steamcmd, "linux32", "steamclient.so")
	dstDir := filepath.Join(config.Force_install_dir, "bin", "lib32")
	dstFile := filepath.Join(dstDir, "steamclient.so")
	backupFile := filepath.Join(dstDir, "steamclient.so.bak")

	// 检查目标文件是否存在
	if _, err := os.Stat(dstFile); err == nil {
		// 如果目标文件存在，先重命名为备份文件
		if err = os.Rename(dstFile, backupFile); err != nil {
			log.Println("重命名文件失败:", err)
			return
		}
		log.Println("已将", dstFile, "重命名为", backupFile)
	}

	// 复制新文件
	if err := copyFile(srcFile, dstFile); err != nil {
		log.Println("复制文件失败:", err)
		return
	}
	log.Println("已将复制到", srcFile, dstFile)

}

// copyFile 实现文件复制
func copyFile(src, dst string) error {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("无法打开源文件: %v", err)
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("无法创建目标文件: %v", err)
	}
	defer dstFile.Close()

	// 复制内容
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("复制内容失败: %v", err)
	}

	return nil
}

// Start 启动世界
func (g *GameLevel2Api) Start(ctx *gin.Context) {
	levelName := ctx.Query("levelName")
	cluster := clusterUtils.GetClusterFromGin(ctx)
	bin := cluster.Bin
	beta := cluster.Beta
	clusterName := cluster.ClusterName

	startBefore(ctx)

	gameService.StartLevel(clusterName, levelName, bin, beta)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "start " + clusterName + " " + levelName + " success",
		Data: nil,
	})
}

// Stop 停止世界
func (g *GameLevel2Api) Stop(ctx *gin.Context) {
	levelName := ctx.Query("levelName")
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	gameService.StopLevel(clusterName, levelName)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "start " + clusterName + " " + levelName + " success",
		Data: nil,
	})
}

// Start 启动世界
func (g *GameLevel2Api) StartAll(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	gameService.StartGame(clusterName)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "start all success",
		Data: nil,
	})
}

// Stop 停止世界
func (g *GameLevel2Api) StopAll(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	gameService.StopGame(clusterName)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "stop all success",
		Data: nil,
	})
}

// GetClusterIni 发送房间配置
func (g *GameLevel2Api) GetClusterIni(ctx *gin.Context) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	clusterIni := ClusterIni{}
	clusterIni.Cluster = homeService.GetClusterIni(clusterName)
	clusterIni.Token = homeService.GetClusterToken(clusterName)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: &clusterIni,
	})
}

// SaveClusterIni 保存房间配置
func (g *GameLevel2Api) SaveClusterIni(ctx *gin.Context) {

	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	clusterIni := ClusterIni{}
	err := ctx.ShouldBind(&clusterIni)
	if err != nil {
		log.Panicln("参数解析错误", err)
	}
	homeService.SaveClusterIni(clusterName, clusterIni.Cluster)
	homeService.SaveClusterToken(clusterName, clusterIni.Token)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: &clusterIni,
	})
}

// SendCommand 发送世界指令
func (g *GameLevel2Api) SendCommand(ctx *gin.Context) {
	var payload struct {
		LevelName string `json:"levelName"`
		Command   string `json:"command"`
	}
	err := ctx.ShouldBind(&payload)
	if err != nil {
		log.Panicln(err)
	}
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	consoleService.SendCommand(clusterName, payload.LevelName, payload.Command)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

// GetAllOnlinePlayers 获取所有在线玩家
func (g *GameLevel2Api) GetAllOnlinePlayers(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	playerList := playerService.GetPlayerList(clusterName, "#ALL_LEVEL")
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: playerList,
	})
}

// GetOnlinePlayers 获取在线玩家
func (g *GameLevel2Api) GetOnlinePlayers(ctx *gin.Context) {
	levelName := ctx.Query("levelName")
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	playerList := playerService.GetPlayerList(clusterName, levelName)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: playerList,
	})
}

// GetAdministrators 获取管理员
func (g *GameLevel2Api) GetAdministrators(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	adminList := playerService.GetDstAdminList(clusterName)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: adminList,
	})
}

// GetWhitelist 获取白名单
func (g *GameLevel2Api) GetWhitelist(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	whitelist := playerService.GetDstWhitelistPlayerList(clusterName)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: whitelist,
	})
}

// GetBlacklist 获取黑名单
func (g *GameLevel2Api) GetBlacklist(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	blacklist := playerService.GetDstBlacklistPlayerList(clusterName)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: blacklist,
	})
}

// SaveBlacklist 保存黑名单
func (g *GameLevel2Api) SaveBlacklist(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	blacklistVO := vo.NewBlacklistVO()
	err := ctx.BindJSON(blacklistVO)
	if err != nil {
		log.Panicln("参数解析错误")
	}

	playerService.SaveBlacklist(clusterName, blacklistVO.Blacklist)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

// SaveWhitelist 保存白名单
func (g *GameLevel2Api) SaveWhitelist(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	whitelist := vo.NewWhitelistVO()
	err := ctx.BindJSON(whitelist)
	if err != nil {
		log.Panicln("参数解析错误")
	}
	playerService.SaveWhitelist(clusterName, whitelist.Whitelist)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

// SaveAdminlist 保存管理员
func (g *GameLevel2Api) SaveAdminlist(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	adminlist := vo.NewAdminListVO()
	err := ctx.BindJSON(adminlist)
	if err != nil {
		log.Panicln("参数解析错误")
	}

	playerService.SaveAdminlist(clusterName, adminlist.AdminList)
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (g *GameLevel2Api) GetLevelList(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: gameLevel2Service.GetLevelList(clusterName),
	})
}

func (g *GameLevel2Api) SaveLevelsList(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName

	var body struct {
		Levels []level.World `json:"levels"`
	}
	err := ctx.ShouldBind(&body)
	if err != nil {
		log.Panicln("参数解析错误", err)
	}

	err = gameLevel2Service.UpdateLevels(clusterName, body.Levels)
	if err != nil {
		log.Panicln("更新世界配置失败", err)
	}

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (g *GameLevel2Api) DeleteLevel(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	levelName := ctx.Query("levelName")

	err := gameLevel2Service.DeleteLevel(clusterName, levelName)
	if err != nil {
		log.Panicln("删除世界失败", err)
	}

	autoCheck.Manager.DeleteAutoCheck(levelName)

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: nil,
	})
}

func (g *GameLevel2Api) CreateNewLevel(ctx *gin.Context) {
	cluster := clusterUtils.GetClusterFromGin(ctx)
	clusterName := cluster.ClusterName
	newLevel := &level.World{}
	err := ctx.ShouldBind(newLevel)
	if err != nil {
		log.Panicln("参数解析错误", err)
	}

	err = gameLevel2Service.CreateLevel(clusterName, newLevel)
	if err != nil {
		log.Panicln("创建世界失败", err)
	}

	autoCheck.Manager.AddAutoCheckTasks(model.AutoCheck{
		ClusterName: clusterName,
		LevelName:   newLevel.LevelName,
		Uuid:        newLevel.Uuid,
		Enable:      0,
		Interval:    10,
		CheckType:   consts.LEVEL_DOWN,
	})

	autoCheck.Manager.AddAutoCheckTasks(model.AutoCheck{
		ClusterName: clusterName,
		LevelName:   newLevel.LevelName,
		Uuid:        newLevel.Uuid,
		Enable:      0,
		Interval:    10,
		CheckType:   consts.LEVEL_MOD,
	})

	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: newLevel,
	})
}

func (g *GameLevel2Api) GetScanUDPPorts(ctx *gin.Context) {
	ports, err := findFreeUDPPorts(10998, 11038)
	if err != nil {
		log.Panicln(err)
	}
	ctx.JSON(http.StatusOK, vo.Response{
		Code: 200,
		Msg:  "success",
		Data: ports,
	})
}

func findFreeUDPPorts(startPort, endPort int) ([]int, error) {
	var freePorts []int

	for port := startPort; port <= endPort; port++ {
		conn, err := net.ListenUDP("udp", &net.UDPAddr{Port: port})
		if err == nil {
			conn.Close()
			freePorts = append(freePorts, port)
		}
	}

	return freePorts, nil
}
