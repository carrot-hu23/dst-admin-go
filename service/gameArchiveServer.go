package service

import (
	"dst-admin-go/constant/dst"
	"dst-admin-go/utils/dstUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type GameArchive struct {
	GameService
	HomeService
}

func (d *GameArchive) GetGameArchive(clusterName string) *vo.GameArchive {

	var wg sync.WaitGroup
	wg.Add(4)

	gameArchie := vo.NewGameArchie()
	basePath := dst.GetClusterBasePath(clusterName)

	// 获取基础信息
	go func() {
		clusterIni := d.GetClusterIni(clusterName)
		gameArchie.ClusterName = clusterIni.ClusterName
		gameArchie.ClusterPassword = clusterIni.ClusterPassword
		gameArchie.GameMod = clusterIni.GameMode
		gameArchie.MaxPlayers = int(clusterIni.MaxPlayers)
		wg.Done()
	}()

	// go func() {
	// 	gameArchie.Players = GetPlayerList()
	// }()

	// 获取mod数量
	go func() {
		masterModoverrides, err := fileUtils.ReadFile(path.Join(basePath, "Master", "modoverrides.lua"))
		if err != nil {
			gameArchie.Mods = 0
		} else {
			gameArchie.Mods = len(dstUtils.WorkshopIds(masterModoverrides))
		}
		wg.Done()
	}()

	// 获取天数和季节
	go func() {
		//metaPath, err := d.FindLatestMetaFile(path.Join(basePath, "Master", "save", "session"))
		//if err != nil {
		//	gameArchie.Meta = ""
		//} else {
		//	meta, err := fileUtils.ReadFile(metaPath)
		//	log.Println("meta path: ", metaPath)
		//	if err != nil {
		//		gameArchie.Meta = ""
		//	} else {
		//		gameArchie.Meta = base64.StdEncoding.EncodeToString([]byte(meta))
		//	}
		//}
		wg.Done()
	}()

	// 获取直连ip
	go func() {
		serverIni := d.GetServerIni(path.Join(basePath, "Master", "server.ini"), true)
		ipv4, err := d.GetPublicIP()
		if err != nil {
			gameArchie.IpConnect = ""
		} else {
			gameArchie.IpConnect = "c_connect(\"" + ipv4 + "\"," + strconv.Itoa(int(serverIni.ServerPort)) + ")"
		}
		wg.Done()
	}()

	wg.Wait()

	return gameArchie
}

func (d *GameArchive) GetPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org/")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}

func (d *GameArchive) getSubPathLevel(rootP, curPath string) int {
	relPath, err := filepath.Rel(rootP, curPath)
	if err != nil {
		// 如果计算相对路径时出错，说明 curPath 不是 rootP 的子目录
		return -1
	}
	// 计算相对路径中 ".." 的数量，即为层数
	return strings.Count(relPath, "..")
}

func (d *GameArchive) FindLatestMetaFile(rootDir string) (string, error) {
	var latestFile string
	var latestModTime time.Time
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".meta" && d.getSubPathLevel(rootDir, path) == 2 {
			if info.ModTime().After(latestModTime) {
				latestFile = path
				latestModTime = info.ModTime()
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return latestFile, nil
}
