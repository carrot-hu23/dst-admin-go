package player

import (
	"dst-admin-go/internal/pkg/utils/fileUtils"
	"dst-admin-go/internal/service/archive"
	"dst-admin-go/internal/service/game"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// PlayerInfo 玩家信息结构体
type PlayerInfo struct {
	Key  string `json:"key"`
	Day  string `json:"day"`
	KuId string `json:"kuId"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type PlayerService struct {
	archive *archive.PathResolver
}

func NewPlayerService(archive *archive.PathResolver) *PlayerService {
	return &PlayerService{
		archive: archive,
	}
}

func (p *PlayerService) GetPlayerList(clusterName string, levelName string, gameProcess game.Process) []PlayerInfo {
	// 处理 #ALL_LEVEL 情况
	queryName := ""
	if levelName == "#ALL_LEVEL" {
		queryName = "Master"
	} else {
		queryName = levelName
	}
	status, _ := gameProcess.Status(clusterName, queryName)
	if !status {
		return make([]PlayerInfo, 0)
	}

	id := strconv.FormatInt(time.Now().Unix(), 10)

	command := ""
	if levelName == "#ALL_LEVEL" {
		levelName = "Master"
		command = "for i, v in ipairs(TheNet:GetClientTable()) do  print(string.format(\\\"player: {[%s] [%d] [%s] [%s] [%s] [%s]} \\\", " + "'" + id + "'" + ",i-1, string.format('%03d', v.playerage), v.userid, v.name, v.prefab)) end"
	} else {
		command = "for i, v in ipairs(AllPlayers) do print(string.format(\\\"player: {[%d] [%d] [%d] [%s] [%s] [%s]} \\\", " + id + ",i,v.components.age:GetAgeInDays(), v.userid, v.name, v.prefab)) end"
	}

	err := gameProcess.Command(clusterName, levelName, command)
	log.Println("clusterName:", clusterName, "levelName:", levelName, "command:", command)
	if err != nil {
		log.Println("Error sending command:", err)
		return make([]PlayerInfo, 0)
	}

	time.Sleep(time.Duration(1) * time.Second)

	// 读取日志
	serverLogPath := p.archive.ServerLogPath(clusterName, levelName)
	dstLogs, err := fileUtils.ReverseRead(serverLogPath, 1000)
	playerInfoList := make([]PlayerInfo, 0)

	for _, line := range dstLogs {
		if strings.Contains(line, id) && strings.Contains(line, "KU") && !strings.Contains(line, "Host") {

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
				if len(result) >= 6 {
					playerInfo := PlayerInfo{Key: result[1], Day: result[2], KuId: result[3], Name: result[4], Role: result[5]}
					playerInfoList = append(playerInfoList, playerInfo)
				}
			}
		}
	}

	// 创建一个map，用于存储不重复的KuId和对应的PlayerInfo对象
	uniquePlayers := make(map[string]PlayerInfo)

	// 遍历players切片
	for _, player := range playerInfoList {
		// 将PlayerInfo对象添加到map中，以KuId作为键
		uniquePlayers[player.KuId] = player
	}

	// 将不重复的PlayerInfo对象从map中提取到新的切片中
	filteredPlayers := make([]PlayerInfo, 0, len(uniquePlayers))
	for _, player := range uniquePlayers {
		filteredPlayers = append(filteredPlayers, player)
	}

	return filteredPlayers
}

func (p *PlayerService) GetPlayerAllList(clusterName string, gameProcess game.Process) []PlayerInfo {
	// 使用 #ALL_LEVEL 调用 GetPlayerList，获取所有玩家
	return p.GetPlayerList(clusterName, "#ALL_LEVEL", gameProcess)
}
