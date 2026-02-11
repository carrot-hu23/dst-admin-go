package service

import (
	"dst-admin-go/constant/screenKey"
	"dst-admin-go/utils/collectionUtils"
	"dst-admin-go/utils/dstUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/utils/shellUtils"
	"dst-admin-go/vo"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type PlayerService struct {
}

func (p *PlayerService) GetAllPlayerList(clusterName string) []vo.PlayerVO {

	if !gameServe.GetLevelStatus(clusterName, "Master") {
		return make([]vo.PlayerVO, 0)
	}

	id := strconv.FormatInt(time.Now().Unix(), 10)

	command := "for i, v in ipairs(TheNet:GetClientTable()) do  print(string.format(\\\"player: {[%s] [%d] [%s] [%s] [%s] [%s]} \\\", " + "'" + id + "'" + ",i-1, string.format('%03d', v.playerage), v.userid, v.name, v.prefab)) end"

	playerCMD := "screen -S \"" + screenKey.Key(clusterName, "Master") + "\" -p 0 -X stuff \"" + command + "\\n\""

	shellUtils.Shell(playerCMD)
	time.Sleep(time.Duration(500) * time.Millisecond)

	// 读取日志
	dstLogs := dstUtils.ReadLevelLog(clusterName, "Master", 100)
	playerVOList := make([]vo.PlayerVO, 0)

	for _, line := range dstLogs {
		if strings.Contains(line, id) && strings.Contains(line, "KU") && !strings.Contains(line, "Host") {
			//str := strings.Split(line, " ")
			//log.Println("players:", str)
			//playerVO := vo.PlayerVO{Key: str[2], Day: str[3], KuId: str[4], Name: str[5], Role: str[6]}
			//playerVOList = append(playerVOList, playerVO)

			log.Println(line)

			// 提取 {} 中的内容
			reCurlyBraces := regexp.MustCompile(`\{([^}]*)\}`)
			curlyBracesMatches := reCurlyBraces.FindStringSubmatch(line)

			if len(curlyBracesMatches) > 1 {
				// curlyBracesMatches[1] 包含 {} 中的内容
				contentInsideCurlyBraces := curlyBracesMatches[1]

				// 提取 [] 中的内容
				reSquareBrackets := regexp.MustCompile(`\[([^\]]*)\]`)
				squareBracketsMatches := reSquareBrackets.FindAllStringSubmatch(contentInsideCurlyBraces, -1)
				var result []string
				for _, match := range squareBracketsMatches {
					// match[1] 包含 [] 中的内容
					contentInsideSquareBrackets := match[1]
					result = append(result, contentInsideSquareBrackets)
				}
				playerVO := vo.PlayerVO{Key: result[1], Day: result[2], KuId: result[3], Name: result[4], Role: result[5]}
				playerVOList = append(playerVOList, playerVO)
			}
		}
	}

	// 创建一个map，用于存储不重复的KuId和对应的PlayerVO对象
	uniquePlayers := make(map[string]vo.PlayerVO)

	// 遍历players切片
	for _, player := range playerVOList {
		// 将PlayerVO对象添加到map中，以KuId作为键
		uniquePlayers[player.KuId] = player
	}

	// 将不重复的PlayerVO对象从map中提取到新的切片中
	filteredPlayers := make([]vo.PlayerVO, 0, len(uniquePlayers))
	for _, player := range uniquePlayers {
		filteredPlayers = append(filteredPlayers, player)
	}

	return filteredPlayers

}

func (p *PlayerService) GetPlayerList(clusterName string, levelName string) []vo.PlayerVO {

	if !gameServe.GetLevelStatus(clusterName, levelName) {
		return make([]vo.PlayerVO, 0)
	}

	id := strconv.FormatInt(time.Now().Unix(), 10)

	command := "for i, v in ipairs(AllPlayers) do print(string.format(\\\"player: {[%d] [%d] [%d] [%s] [%s] [%s]} \\\", " + id + ",i,v.components.age:GetAgeInDays(), v.userid, v.name, v.prefab)) end"
	// command := "for i, v in ipairs(TheNet:GetClientTable()) do  print(string.format(\\\"%s %d %s %s %s %s \\\", " + "'" + id + "'" + ",i-1, string.format('%03d', v.playerage), v.userid, v.name, v.prefab)) end"

	playerCMD := "screen -S \"" + screenKey.Key(clusterName, levelName) + "\" -p 0 -X stuff \"" + command + "\\n\""

	shellUtils.Shell(playerCMD)
	time.Sleep(time.Duration(500) * time.Millisecond)

	// 读取日志
	dstLogs := dstUtils.ReadLevelLog(clusterName, levelName, 100)
	playerVOList := make([]vo.PlayerVO, 0)

	for _, line := range dstLogs {
		if strings.Contains(line, id) && strings.Contains(line, "KU") && !strings.Contains(line, "Host") {
			//str := strings.Split(line, " ")
			//log.Println("players:", str)
			//playerVO := vo.PlayerVO{Key: str[2], Day: str[3], KuId: str[4], Name: str[5], Role: str[6]}
			//playerVOList = append(playerVOList, playerVO)

			log.Println(line)

			// 提取 {} 中的内容
			reCurlyBraces := regexp.MustCompile(`\{([^}]*)\}`)
			curlyBracesMatches := reCurlyBraces.FindStringSubmatch(line)

			if len(curlyBracesMatches) > 1 {
				// curlyBracesMatches[1] 包含 {} 中的内容
				contentInsideCurlyBraces := curlyBracesMatches[1]

				// 提取 [] 中的内容
				reSquareBrackets := regexp.MustCompile(`\[([^\]]*)\]`)
				squareBracketsMatches := reSquareBrackets.FindAllStringSubmatch(contentInsideCurlyBraces, -1)
				var result []string
				for _, match := range squareBracketsMatches {
					// match[1] 包含 [] 中的内容
					contentInsideSquareBrackets := match[1]
					result = append(result, contentInsideSquareBrackets)
				}
				playerVO := vo.PlayerVO{Key: result[1], Day: result[2], KuId: result[3], Name: result[4], Role: result[5]}
				playerVOList = append(playerVOList, playerVO)
			}
		}
	}

	// 创建一个map，用于存储不重复的KuId和对应的PlayerVO对象
	uniquePlayers := make(map[string]vo.PlayerVO)

	// 遍历players切片
	for _, player := range playerVOList {
		// 将PlayerVO对象添加到map中，以KuId作为键
		uniquePlayers[player.KuId] = player
	}

	// 将不重复的PlayerVO对象从map中提取到新的切片中
	filteredPlayers := make([]vo.PlayerVO, 0, len(uniquePlayers))
	for _, player := range uniquePlayers {
		filteredPlayers = append(filteredPlayers, player)
	}

	return filteredPlayers

}

func (p *PlayerService) GetDstAdminList(clusterName string) (str []string) {
	path := dstUtils.GetAdminlistPath(clusterName)
	if !fileUtils.Exists(path) {
		log.Println("路径不存在", path)
		return
	}
	str, err := fileUtils.ReadLnFile(path)
	log.Println("str:", str)
	if err != nil {
		panic("read dstUtils2 adminlist.txt error: \n" + err.Error())
	}
	return
}

func (p *PlayerService) GetDstBlacklistPlayerList(clusterName string) (str []string) {
	path := dstUtils.GetBlocklistPath(clusterName)
	if !fileUtils.Exists(path) {
		log.Println("路径不存在", path)
		return
	}
	str, err := fileUtils.ReadLnFile(path)
	log.Println("str:", str)
	if err != nil {
		panic("read dstUtils2 blocklist.txt error: \n" + err.Error())
	}
	return
}

func (p *PlayerService) GetDstWhitelistPlayerList(clusterName string) (str []string) {
	path := dstUtils.GetWhitelistPath(clusterName)
	if !fileUtils.Exists(path) {
		log.Println("路径不存在", path)
		return
	}
	str, err := fileUtils.ReadLnFile(path)
	log.Println("str:", str)
	if err != nil {
		panic("read dstUtils2 whitelist.txt error: \n" + err.Error())
	}
	return
}

func (p *PlayerService) SaveDstAdminList(clusterName string, adminlist []string) {

	path := dstUtils.GetAdminlistPath(clusterName)

	err := fileUtils.CreateFileIfNotExists(path)
	if err != nil {
		panic("create dstUtils2 blacklist.txt error: \n" + err.Error())
	}
	lnFile, err := fileUtils.ReadLnFile(path)
	set := collectionUtils.ToSet(append(lnFile, adminlist...))

	err = fileUtils.WriterLnFile(path, set)
	if err != nil {
		panic("write dstUtils2 blacklist.txt error: \n" + err.Error())
	}
}

func (p *PlayerService) SaveDstBlacklistPlayerList(clusterName string, blacklist []string) {

	path := dstUtils.GetBlocklistPath(clusterName)

	err := fileUtils.CreateFileIfNotExists(path)
	if err != nil {
		panic("create dstUtils2 blacklist.txt error: \n" + err.Error())
	}
	lnFile, err := fileUtils.ReadLnFile(path)
	set := collectionUtils.ToSet(append(lnFile, blacklist...))

	err = fileUtils.WriterLnFile(path, set)
	if err != nil {
		panic("write dstUtils2 blacklist.txt error: \n" + err.Error())
	}
}

func (p *PlayerService) DeleteDstBlacklistPlayerList(clusterName string, blacklist []string) {

	path := dstUtils.GetBlocklistPath(clusterName)
	err := fileUtils.CreateFileIfNotExists(path)
	if err != nil {
		panic("create dstUtils2 adminlist.txt error: \n" + err.Error())
	}
	lnFile, err := fileUtils.ReadLnFile(path)
	var result []string
	for i := range lnFile {
		isFind := false
		for j := range blacklist {
			if lnFile[i] == blacklist[j] {
				isFind = true
				break
			}
		}
		if !isFind {
			result = append(result, lnFile[i])
		}
	}

	err = fileUtils.WriterLnFile(path, result)
	if err != nil {
		panic("write dstUtils2 adminlist.txt error: \n" + err.Error())
	}
}

func (p *PlayerService) DeleteDstAdminListPlayerList(clusterName string, adminlist []string) {

	path := dstUtils.GetAdminlistPath(clusterName)
	err := fileUtils.CreateFileIfNotExists(path)
	if err != nil {
		panic("create dstUtils2 adminlist.txt error: \n" + err.Error())
	}
	lnFile, err := fileUtils.ReadLnFile(path)
	var result []string
	for i := range lnFile {
		isFind := false
		for j := range adminlist {
			if lnFile[i] == adminlist[j] {
				isFind = true
				break
			}
		}
		if !isFind {
			result = append(result, lnFile[i])
		}
	}

	err = fileUtils.WriterLnFile(path, result)
	if err != nil {
		panic("write dstUtils2 adminlist.txt error: \n" + err.Error())
	}

}

func (p *PlayerService) SaveBlacklist(clusterName string, list []string) {

	path := dstUtils.GetBlacklistPath(clusterName)

	err := fileUtils.CreateFileIfNotExists(path)
	if err != nil {
		panic("create blacklist.txt error: \n" + err.Error())
	}

	err = fileUtils.WriterLnFile(path, list)
	if err != nil {
		panic("write dstUtils2 blacklist.txt error: \n" + err.Error())
	}
}

func (p *PlayerService) SaveWhitelist(clusterName string, list []string) {

	path := dstUtils.GetWhitelistPath(clusterName)

	err := fileUtils.CreateFileIfNotExists(path)
	if err != nil {
		panic("create blacklist.txt error: \n" + err.Error())
	}

	err = fileUtils.WriterLnFile(path, list)
	if err != nil {
		panic("write dstUtils2 blacklist.txt error: \n" + err.Error())
	}
}

func (p *PlayerService) SaveAdminlist(clusterName string, list []string) {

	path := dstUtils.GetAdminlistPath(clusterName)

	err := fileUtils.CreateFileIfNotExists(path)
	if err != nil {
		panic("create blacklist.txt error: \n" + err.Error())
	}

	err = fileUtils.WriterLnFile(path, list)
	if err != nil {
		panic("write dstUtils2 blacklist.txt error: \n" + err.Error())
	}
}
