package service

import (
	"dst-admin-go/constant"
	"dst-admin-go/utils/dstConfigUtils"
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

func (p *PlayerService) GetPlayerList() []vo.PlayerVO {
	id := strconv.FormatInt(time.Now().Unix(), 10)

	command := "for i, v in ipairs(TheNet:GetClientTable()) do  print(string.format(\\\"%s %d %s %s %s %s \\\", " + "'" + id + "'" + ",i-1, string.format('%03d', v.playerage), v.userid, v.name, v.prefab)) end"
	clsuerName := dstConfigUtils.GetDstConfig().Cluster
	screenKey := getscreenKey(clsuerName, "Master")

	playerCMD := "screen -S \"" + screenKey + "\" -p 0 -X stuff \"" + command + "\\n\""

	shellUtils.Shell(playerCMD)

	time.Sleep(time.Duration(1) * time.Second)

	// 读取日志
	dstLogs := ReadDstMasterLog(100)
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

func (p *PlayerService) GetDstAdminList() (str []string) {
	path := constant.GET_DST_ADMIN_LIST_PATH()
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

func (p *PlayerService) GetDstBlcaklistPlayerList() (str []string) {
	path := constant.GET_DST_BLOCKLIST_PATH()
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

func (p *PlayerService) SaveDstAdminList(adminlist []string) {

	path := constant.GET_DST_ADMIN_LIST_PATH()

	err := fileUtils.WriterLnFile(path, adminlist)
	if err != nil {
		panic("write dst adminlist.txt error: \n" + err.Error())
	}
}

func (p *PlayerService) SaveDstBlacklistPlayerList(blacklist []string) {

	path := constant.GET_DST_BLOCKLIST_PATH()

	err := fileUtils.WriterLnFile(path, blacklist)
	if err != nil {
		panic("write dst adminlist.txt error: \n" + err.Error())
	}
}
