package service

import (
	"dst-admin-go/constant/dst"
	"dst-admin-go/constant/screenKey"
	"dst-admin-go/utils/dstUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/utils/shellUtils"
	"dst-admin-go/vo"
	"log"
	"strconv"
	"strings"
	"time"
)

type PlayerService struct {
}

func (p *PlayerService) GetPlayerList(clusterName string) []vo.PlayerVO {
	id := strconv.FormatInt(time.Now().Unix(), 10)

	command := "for i, v in ipairs(TheNet:GetClientTable()) do  print(string.format(\\\"%s %d %s %s %s %s \\\", " + "'" + id + "'" + ",i-1, string.format('%03d', v.playerage), v.userid, v.name, v.prefab)) end"

	playerCMD := "screen -S \"" + screenKey.Key(clusterName, "Master") + "\" -p 0 -X stuff \"" + command + "\\n\""

	shellUtils.Shell(playerCMD)

	time.Sleep(time.Duration(1) * time.Second)

	// TODO 如果只启动了洞穴，应该去读取洞穴的日志

	// 读取日志
	dstLogs := dstUtils.ReadMasterLog(clusterName, 100)
	var playerVOList []vo.PlayerVO

	for _, line := range dstLogs {
		if strings.Contains(line, id) && strings.Contains(line, "KU") && !strings.Contains(line, "Host") {
			str := strings.Split(line, " ")
			log.Println("players:", str)
			playerVO := vo.PlayerVO{Key: str[2], Day: str[3], KuId: str[4], Name: str[5], Role: str[6]}
			playerVOList = append(playerVOList, playerVO)
		}
	}

	return playerVOList
}

func (p *PlayerService) GetDstAdminList(clusterName string) (str []string) {
	path := dst.GetAdminlistPath(clusterName)
	if !fileUtils.Exists(path) {
		log.Println("路径不存在", path)
		return
	}
	str, err := fileUtils.ReadLnFile(path)
	log.Println("str:", str)
	if err != nil {
		panic("read dst adminlist.txt error: \n" + err.Error())
	}
	return
}

func (p *PlayerService) GetDstBlcaklistPlayerList(clusterName string) (str []string) {
	path := dst.GetBlocklistPath(clusterName)
	if !fileUtils.Exists(path) {
		log.Println("路径不存在", path)
		return
	}
	str, err := fileUtils.ReadLnFile(path)
	log.Println("str:", str)
	if err != nil {
		panic("read dst blocklist.txt error: \n" + err.Error())
	}
	return
}

func (p *PlayerService) SaveDstAdminList(clusterName string, adminlist []string) {

	path := dst.GetAdminlistPath(clusterName)

	err := fileUtils.WriterLnFile(path, adminlist)
	if err != nil {
		panic("write dst adminlist.txt error: \n" + err.Error())
	}
}

func (p *PlayerService) SaveDstBlacklistPlayerList(clusterName string, blacklist []string) {

	path := dst.GetBlocklistPath(clusterName)

	err := fileUtils.WriterLnFile(path, blacklist)
	if err != nil {
		panic("write dst adminlist.txt error: \n" + err.Error())
	}
}
