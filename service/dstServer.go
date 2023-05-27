package service

import (
	"dst-admin-go/constant"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/vo"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

func GetCurrGameArchive() *vo.GameArchive {

	var wg sync.WaitGroup
	wg.Add(4)

	gameArchie := vo.NewGameArchie()
	basePath := constant.GET_DST_USER_GAME_CONFG_PATH()

	// 获取基础信息
	go func() {
		clusterIni := ReadClusterIniFile()
		gameArchie.ClusterName = clusterIni.ClusterName
		gameArchie.ClusterPassword = clusterIni.ClusterPassword
		gameArchie.GameMod = clusterIni.GameMode
		gameArchie.MaxPlayers = int(clusterIni.MaxPlayers)
		wg.Done()
	}()

	go func() {
		gameArchie.Players = GetPlayerList()
	}()

	// 获取mod数量
	go func() {
		if getCavesStatus() {
			cavesModoverrides, err := fileUtils.ReadFile(path.Join(basePath, "Caves", "modoverrides.lua"))
			if err != nil {
				gameArchie.Mods = 0
			} else {
				gameArchie.Mods = len(WorkshopIds(cavesModoverrides))
			}

		} else {
			masterModoverrides, err := fileUtils.ReadFile(path.Join(basePath, "Master", "modoverrides.lua"))
			if err != nil {
				gameArchie.Mods = 0
			} else {
				gameArchie.Mods = len(WorkshopIds(masterModoverrides))
			}
		}

		wg.Done()
	}()

	// 获取天数和季节
	go func() {
		metaPath, err := FindLatestMetaFile(path.Join(basePath, "Master", "save", "session"))
		if err != nil {
			gameArchie.Meta = ""
		} else {
			meta, err := fileUtils.ReadFile(metaPath)
			log.Println("meta path: ", metaPath)
			if err != nil {
				gameArchie.Meta = ""
			} else {
				gameArchie.Meta = base64.StdEncoding.EncodeToString([]byte(meta))
			}
		}
		wg.Done()
	}()

	// 获取直连ip
	go func() {
		serverIni := ReadServerIniFile(path.Join(basePath, "Master", "server.ini"), true)
		ipv4, err := GetPublicIP()
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

// func FindLatestMetaFile(basePath string) (string, error) {

// 	metaPath := path.Join(basePath, "Master", "save", "session")

// 	var newestMetaPath string
// 	var newestModTime time.Time

// 	err := filepath.Walk(metaPath, func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}
// 		if !info.IsDir() && filepath.Ext(path) == ".meta" {
// 			modTime := info.ModTime()
// 			if modTime.After(newestModTime) {
// 				newestMetaPath = path
// 				newestModTime = modTime
// 			}
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		return "", err
// 	}

// 	if newestMetaPath == "" {
// 		return "", fmt.Errorf("No .meta file found")
// 	} else {
// 		return newestMetaPath, nil
// 	}
// }

func GetPublicIP() (string, error) {
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

func FindLatestMetaFile(rootDir string) (string, error) {
	var latestFile string
	var latestModTime time.Time
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".meta" {
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
