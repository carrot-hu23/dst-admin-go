package service

import (
	"dst-admin-go/constant/consts"
	"dst-admin-go/utils/dstUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/utils/luaUtils"
	"dst-admin-go/vo"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

type GameArchive struct {
	GameService
	HomeService
	PlayerService
}

func (d *GameArchive) GetGameArchive(clusterName string) *vo.GameArchive {

	var wg sync.WaitGroup
	wg.Add(6)

	gameArchie := vo.NewGameArchie()
	basePath := dstUtils.GetClusterBasePath(clusterName)

	// 获取基础信息
	go func() {
		clusterIni := d.GetClusterIni(clusterName)
		gameArchie.ClusterName = clusterIni.ClusterName
		gameArchie.ClusterDescription = clusterIni.ClusterDescription
		gameArchie.ClusterPassword = clusterIni.ClusterPassword
		gameArchie.GameMod = clusterIni.GameMode
		gameArchie.MaxPlayers = int(clusterIni.MaxPlayers)
		wg.Done()
	}()

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
		defer func() {
			wg.Done()
			if r := recover(); r != nil {
			}
		}()
		gameArchie.Meta = d.Snapshoot(clusterName)
	}()

	// 获取直连ip
	go func() {
		defer func() {
			wg.Done()
			if r := recover(); r != nil {

			}
		}()
		clusterIni := d.GetClusterIni(clusterName)
		password := clusterIni.ClusterPassword
		serverIni := d.GetServerIni(path.Join(basePath, "Master", "server.ini"), true)
		ipv4, err := d.GetPublicIP()
		if err != nil {
			gameArchie.IpConnect = ""
		} else {
			// c_connect("IP address", port, "password")
			if password != "" {
				gameArchie.IpConnect = "c_connect(\"" + ipv4 + "\"," + strconv.Itoa(int(serverIni.ServerPort)) + ",\"" + password + "\"" + ")"
			} else {
				gameArchie.IpConnect = "c_connect(\"" + ipv4 + "\"," + strconv.Itoa(int(serverIni.ServerPort)) + ")"
			}
		}
		gameArchie.Port = serverIni.ServerPort
		gameArchie.Ip = ipv4

	}()

	go func() {
		defer func() {
			wg.Done()
			if r := recover(); r != nil {

			}
		}()
		localVersion := d.GetLocalDstVersion(clusterName)
		version := d.GetLastDstVersion()

		gameArchie.Version = localVersion
		gameArchie.LastVersion = version
	}()

	go func() {
		defer func() {
			wg.Done()
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		// TODO 默认取Master世界人数
		gameArchie.Players = d.GetPlayerList(clusterName, "Master")
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

func findLatestMetaFile(directory string) (string, error) {
	// 检查指定目录是否存在
	_, err := os.Stat(directory)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("目录不存在：%s", directory)
	}

	// 获取指定目录下一级的所有子目录
	subdirs, err := ioutil.ReadDir(directory)
	if err != nil {
		return "", fmt.Errorf("读取目录失败：%s", err)
	}

	// 用于存储最新的.meta文件路径和其修改时间
	var latestMetaFile string
	var latestMetaFileTime time.Time

	for _, subdir := range subdirs {
		// 检查子目录是否是目录
		if subdir.IsDir() {
			subdirPath := filepath.Join(directory, subdir.Name())

			// 获取子目录下的所有文件
			files, err := ioutil.ReadDir(subdirPath)
			if err != nil {
				return "", fmt.Errorf("读取子目录失败：%s", err)
			}

			for _, file := range files {
				// 检查文件是否是.meta文件
				if !file.IsDir() && filepath.Ext(file.Name()) == ".meta" {
					// 获取文件的修改时间
					modifiedTime := file.ModTime()

					// 如果找到的文件的修改时间比当前最新的.meta文件的修改时间更晚，则更新最新的.meta文件路径和修改时间
					if modifiedTime.After(latestMetaFileTime) {
						latestMetaFile = filepath.Join(subdirPath, file.Name())
						latestMetaFileTime = modifiedTime
					}
				}
			}
		}
	}

	if latestMetaFile == "" {
		return "", fmt.Errorf("未找到.meta文件")
	}

	return latestMetaFile, nil
}

func (d *GameArchive) Snapshoot(clusterName string) vo.Meta {
	sessionPath := filepath.Join(consts.KleiDstPath, clusterName, "Master", "save", "session")
	p, err := findLatestMetaFile(sessionPath)
	if err != nil {
		fmt.Println("查找meta文件失败", err)
		return vo.Meta{}
	}
	content, err := fileUtils.ReadFile(p)
	if err != nil {
		fmt.Println("读取meta文件失败", err)
		return vo.Meta{}
	}
	var data vo.Meta
	err = luaUtils.LuaTable2Struct(content[:len(content)-1], reflect.ValueOf(&data).Elem())
	if err != nil {
		fmt.Println("解析meta文件失败", err)
		return vo.Meta{}
	}
	return data
}
