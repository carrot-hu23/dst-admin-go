package service

import (
	"context"
	"dst-admin-go/constant/screenKey"
	"dst-admin-go/utils/collectionUtils"
	"dst-admin-go/utils/dstUtils"
	"dst-admin-go/utils/fileUtils"
	"dst-admin-go/utils/shellUtils"
	"dst-admin-go/vo"
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hpcloud/tail"
)

type result struct {
	value interface{}
	err   error
}

type entry struct {
	res        result
	ready      chan struct{} // closed when res is ready
	expiration time.Time     // expiration time
}

type Func func(key string) (interface{}, error)

func New(f Func) *Memo {
	return &Memo{f: f, cache: make(map[string]*entry)}
}

type Memo struct {
	f     Func
	mu    sync.Mutex // guards cache
	cache map[string]*entry
}

func (memo *Memo) Get(key string) (value interface{}, err error) {
	memo.mu.Lock()
	e := memo.cache[key]
	if e == nil || e.expiration.Before(time.Now()) {
		// This is the first request for this key or it has expired.
		// This goroutine becomes responsible for computing
		// the value and broadcasting the ready condition.
		e = &entry{
			ready:      make(chan struct{}),
			expiration: time.Now().Add(5 * time.Second), // set expiration time to 1 minute from now
		}
		memo.cache[key] = e
		memo.mu.Unlock()

		e.res.value, e.res.err = memo.f(key)
		close(e.ready) // broadcast ready condition
	} else {
		// This is a repeat request for this key and it has not expired.
		memo.mu.Unlock()

		<-e.ready // wait for ready condition
	}
	return e.res.value, e.res.err
}

var memo *Memo

func init() {
	memo = New(func(key string) (interface{}, error) {

		if isWindows() {
			split := strings.Split(key, ":")
			if len(split) != 2 {
				return []vo.PlayerVO{}, nil
			}
			clusterName := split[0]
			levelName := split[1]

			if levelName != "#ALL_LEVEL" {
				if !gameServe.GetLevelStatus(clusterName, levelName) {
					return []vo.PlayerVO{}, nil
				}
			}

			id := strconv.FormatInt(time.Now().Unix(), 10)

			if levelName == "#ALL_LEVEL" {
				levelName = "Master"
			}
			command := "for i, v in ipairs(TheNet:GetClientTable()) do  print(string.format(\"%s %d %s %s %s %s \", " + "'" + id + "'" + ",i-1, string.format('%03d', v.playerage), v.userid, v.name, v.prefab)) end"

			clusterContainer.Send(clusterName, levelName, command)
			time.Sleep(time.Duration(1) * time.Second)
			// 读取日志
			dstLogs := dstUtils.ReadLevelLog(clusterName, levelName, 100)
			playerVOList := make([]vo.PlayerVO, 0)

			for _, line := range dstLogs {
				if strings.Contains(line, id) && strings.Contains(line, "KU") && !strings.Contains(line, "Host") {
					str := strings.Split(line, " ")
					log.Println("players:", str)
					playerVO := vo.PlayerVO{Key: str[2], Day: str[3], KuId: str[4], Name: str[5], Role: str[6]}
					playerVOList = append(playerVOList, playerVO)
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

			return filteredPlayers, nil
		}

		// #################################################################

		split := strings.Split(key, ":")
		log.Println("-----------", key)
		if len(split) != 2 {
			return []vo.PlayerVO{}, nil
		}
		clusterName := split[0]
		levelName := split[1]
		if levelName != "#ALL_LEVEL" {
			if !gameServe.GetLevelStatus(clusterName, levelName) {
				return []vo.PlayerVO{}, nil
			}
		}

		id := strconv.FormatInt(time.Now().Unix(), 10)
		command := ""
		if levelName == "#ALL_LEVEL" {
			levelName = "Master"
			command = "for i, v in ipairs(TheNet:GetClientTable()) do  print(string.format(\\\"player: {[%s] [%d] [%s] [%s] [%s] [%s]} \\\", " + "'" + id + "'" + ",i-1, string.format('%03d', v.playerage), v.userid, v.name, v.prefab)) end"
		} else {
			command = "for i, v in ipairs(AllPlayers) do print(string.format(\\\"player: {[%d] [%d] [%d] [%s] [%s] [%s]} \\\", " + id + ",i,v.components.age:GetAgeInDays(), v.userid, v.name, v.prefab)) end"
		}
		// command := "for i, v in ipairs(TheNet:GetClientTable()) do  print(string.format(\\\"%s %d %s %s %s %s \\\", " + "'" + id + "'" + ",i-1, string.format('%03d', v.playerage), v.userid, v.name, v.prefab)) end"
		// command := "for i, v in ipairs(AllPlayers) do  print(string.format(\\\"%s %d %s %s %s %s \\\", " + "'" + id + "'" + ",i-1, string.format('%03d', v.playerage), v.userid, v.name, v.prefab)) end"

		playerCMD := "screen -S \"" + screenKey.Key(clusterName, levelName) + "\" -p 0 -X stuff \"" + command + "\\n\""
		log.Println("playerCMD", playerCMD)
		shellUtils.Shell(playerCMD)

		time.Sleep(time.Duration(1) * time.Second)

		// TODO 如果只启动了洞穴，应该去读取洞穴的日志

		// 读取日志
		dstLogs := dstUtils.ReadLevelLog(clusterName, levelName, 150)
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
		return filteredPlayers, nil
	})
}

type PlayerService struct {
}

func (p *PlayerService) GetPlayerList(clusterName string, levelName string) []vo.PlayerVO {
	originalLevelName := levelName
	if levelName == "#ALL_LEVEL" {
		levelName = "Master"
	}

	status := gameServe.GetLevelStatus(clusterName, levelName)
	if !status {
		return make([]vo.PlayerVO, 0)
	}

	// 步骤 1: 准备唯一ID和结束信号
	id := fmt.Sprintf("req-%d", time.Now().UnixNano())
	endSignal := fmt.Sprintf("end-signal-for-%s", id)

	// 步骤 2: 修改Command生成逻辑
	command := ""
	endCommand := fmt.Sprintf("; print(\\\"%s\\\")", endSignal)

	if originalLevelName == "#ALL_LEVEL" {
		levelName = "Master"
		baseCommand := "for i, v in ipairs(TheNet:GetClientTable()) do  print(string.format(\\\"player: {[%s] [%d] [%s] [%s] [%s] [%s]} \\\", " + "'" + id + "'" + ",i-1, string.format('%03d', v.playerage), v.userid, v.name, v.prefab)) end"
		//baseCommand := fmt.Sprintf("for i, v in ipairs(TheNet:GetClientTable()) do print(string.format(\\\"player: {[%s] [%d] [%s] [%s] [%s] [%s]} \\\", '%s', i-1, string.format('%%03d', v.playerage), v.userid, v.name, v.prefab)) end", id, id)
		command = baseCommand + endCommand
	} else {
		//baseCommand := fmt.Sprintf("for i, v in ipairs(AllPlayers) do print(string.format(\\\"player: {[%s] [%d] [%d] [%s] [%s] [%s]} \\\", '%s', i,v.components.age:GetAgeInDays(), v.userid, v.name, v.prefab)) end", id, id)
		baseCommand := "for i, v in ipairs(AllPlayers) do print(string.format(\\\"player: {[%s] [%d] [%d] [%s] [%s] [%s]} \\\", " + "'" + id + "'" + ",i,v.components.age:GetAgeInDays(), v.userid, v.name, v.prefab)) end"
		command = baseCommand + endCommand
	}

	// 步骤 3: 替换核心逻辑
	// -------------------- 新逻辑开始 --------------------
	//logFilePath := dstUtils.GetLevelLogPath(clusterName, levelName) // 假设你有这么一个函数获取日志路径
	logFilePath := filepath.Join(dstUtils.GetClusterBasePath(clusterName), levelName, "server_log.txt")
	t, err := tail.TailFile(logFilePath, tail.Config{
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: true,
		Follow:    true,
	})
	if err != nil {
		log.Printf("错误：无法追踪日志文件 %s: %v", logFilePath, err)
		return make([]vo.PlayerVO, 0)
	}
	defer t.Stop()

	// 设置一个10秒的超时作为安全网
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 发送命令 (你原有的逻辑)
	if isWindows() {
		clusterContainer.Send(clusterName, levelName, command)
	} else {
		playerCMD := "screen -S \"" + screenKey.Key(clusterName, levelName) + "\" -p 0 -X stuff \"" + command + "\\n\""
		shellUtils.Shell(playerCMD)
	}

	log.Printf("[命令已发送] 开始监听日志中的 ID: %s", id)
	uniquePlayers := make(map[string]vo.PlayerVO) // 使用map实现自动去重

	for {
		select {
		case line, ok := <-t.Lines:
			if !ok {
				log.Printf("[日志追踪] tail channel 已关闭，即将退出。")
				goto EndLoop
			}

			// 如果是服务器的命令回显日志 或者是主机，则直接跳过，等待下一行
			if strings.Contains(line.Text, "RemoteCommandInput") || strings.Contains(line.Text, "Host") {
				continue
			}

			// 检查是否是结束信号
			if strings.Contains(line.Text, endSignal) {
				log.Printf("[信号捕获] 成功匹配到结束信号 for ID: %s", id)
				goto EndLoop
			}

			// 检查是否是数据信号
			if strings.Contains(line.Text, id) && strings.Contains(line.Text, "KU") {
				// 正则解析逻辑
				reCurlyBraces := regexp.MustCompile(`\{([^}]*)\}`)
				curlyBracesMatches := reCurlyBraces.FindStringSubmatch(line.Text)
				if len(curlyBracesMatches) > 1 {
					contentInsideCurlyBraces := curlyBracesMatches[1]
					reSquareBrackets := regexp.MustCompile(`\[([^\]]*)\]`)
					squareBracketsMatches := reSquareBrackets.FindAllStringSubmatch(contentInsideCurlyBraces, -1)
					var result []string
					for _, match := range squareBracketsMatches {
						result = append(result, match[1])
					}
					if len(result) >= 6 {
						playerVO := vo.PlayerVO{Key: result[1], Day: result[2], KuId: result[3], Name: result[4], Role: result[5]}
						uniquePlayers[playerVO.KuId] = playerVO
					}
				}
			}
		case <-ctx.Done():
			log.Printf("警告：在超时前未收到结束信号 for ID: %s。可能命令执行失败或服务器繁忙。", id)
			goto EndLoop
		}
	}

EndLoop:
	// 7. 将去重后的map转换为slice返回
	log.Printf("监控结束 for ID: %s。共捕获到 %d 名不重复的玩家。", id, len(uniquePlayers))
	filteredPlayers := make([]vo.PlayerVO, 0, len(uniquePlayers))
	for _, player := range uniquePlayers {
		filteredPlayers = append(filteredPlayers, player)
	}

	// --- 8. (新增的关键步骤) 对最终的 slice 按 KuId 进行排序 ---
	sort.Slice(filteredPlayers, func(i, j int) bool {
		// 返回 true 表示 i 应该排在 j 前面
		// 这里我们按 KuId 的字母顺序升序排列
		return filteredPlayers[i].KuId < filteredPlayers[j].KuId
	})
	// ----------------------------------------------------

	// 9. 返回一个干净、去重且排序好的列表
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
