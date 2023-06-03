package collect

//import (
//	"dst-admin-go/chatgpt"
//	"dst-admin-go/model"
//	"dst-admin-go/service"
//	"fmt"
//	"log"
//	"path/filepath"
//	"strings"
//	"time"
//
//	"github.com/hpcloud/tail"
//	"golang.org/x/text/encoding/simplifiedchinese"
//)
//
//func CollectChatLog(text string) {
//
//	defer func() {
//		if err := recover(); err != nil {
//			log.Println("玩家行为日志解析异常:", err)
//		}
//	}()
//
//	//[00:00:55]: [Join Announcement] 猜猜我是谁
//	if strings.Contains(text, "[Join Announcement]") {
//		parseJoin(text)
//	}
//	//[00:02:28]: [Leave Announcement] 猜猜我是谁
//	if strings.Contains(text, "[Leave Announcement]") {
//		parseLeave(text)
//	}
//	//[00:02:17]: [Death Announcement] 猜猜我是谁 死于： 采摘的红蘑菇。她变成了可怕的鬼魂！
//	if strings.Contains(text, "[Death Announcement]") {
//		parseDeath(text)
//	}
//	//[00:02:37]: [Resurrect Announcement] 猜猜我是谁 复活自： TMIP 控制台.
//	if strings.Contains(text, "[Resurrect Announcement]") {
//		parseResurrect(text)
//	}
//	//[00:03:16]: [Say] (KU_Mt-zrX8K) 猜猜我是谁: 你好啊
//	if strings.Contains(text, "[Say]") {
//		parseSay(text)
//	}
//}
//
//func parseSay(text string) {
//	fmt.Println(text)
//
//	arr := strings.Split(text, " ")
//	temp := strings.Replace(arr[0], " ", "", -1)
//	time := temp[:len(temp)-1]
//	action := arr[1]
//	kuId := arr[2]
//	kuId = kuId[1 : len(kuId)-1]
//	name := arr[3]
//	name = name[:len(name)-1]
//	rest := ""
//	for i := 4; i <= len(arr)-1; i++ {
//		rest += arr[i] + " "
//	}
//	actionDesc := rest
//
//	spawn := getSpawnRole(name)
//	connect := getConnectInfo(name)
//
//	playerLog := model.PlayerLog{
//		Name:       name,
//		Role:       spawn.Role,
//		Action:     action,
//		ActionDesc: actionDesc,
//		Time:       time,
//		Ip:         connect.Ip,
//		KuId:       kuId,
//		SteamId:    connect.SteamId,
//	}
//	//fmt.Println("time", time, "action:", action, "name:", "kuId:", kuId, name, "op:", actionDesc)
//	//获取最近的一条spwan记录和newComing
//	//playerLog := model.PlayerLog{Name: name, Role: spawn.Role, KuId: kuId, Action: action, ActionDesc: actionDesc, Time: time}
//
//	model.DB.Create(&playerLog)
//
//	if strings.Contains(text, model.Config.Flag) {
//		arr := strings.Split(text, model.Config.Flag)
//		s := arr[1]
//		chatgpt.ChatGpt(kuId, s, service.SentBroadcast)
//	}
//
//}
//
//func parseResurrect(text string) {
//	parseDeath(text)
//}
//
//func parseDeath(text string) {
//	fmt.Println(text)
//	arr := strings.Split(text, " ")
//
//	temp := strings.Replace(arr[0], " ", "", -1)
//	time := temp[:len(temp)-1]
//	action := arr[1] + arr[2]
//	name := strings.Replace(arr[3], "\n", "", -1)
//
//	rest := ""
//	for i := 4; i <= len(arr)-1; i++ {
//		rest += arr[i] + " "
//	}
//	actionDesc := rest
//
//	//获取最近的一条spwan记录和newComing
//	spawn := getSpawnRole(name)
//	connect := getConnectInfo(name)
//	fmt.Println(connect)
//
//	playerLog := model.PlayerLog{
//		Name:       name,
//		Role:       spawn.Role,
//		Action:     action,
//		ActionDesc: actionDesc,
//		Time:       time,
//		Ip:         connect.Ip,
//		KuId:       connect.KuId,
//		SteamId:    connect.SteamId,
//	}
//
//	model.DB.Create(&playerLog)
//
//}
//
//func parseLeave(text string) {
//	parseJoin(text)
//}
//
//func parseJoin(text string) {
//	fmt.Println(text)
//	arr := strings.Split(text, " ")
//	temp := strings.Replace(arr[0], " ", "", -1)
//	time := temp[:len(temp)-1]
//	action := arr[1] + arr[2]
//	name := arr[3]
//
//	spawn := getSpawnRole(name)
//	connect := getConnectInfo(name)
//
//	playerLog := model.PlayerLog{
//		Name:    name,
//		Role:    spawn.Role,
//		Action:  action,
//		Time:    time,
//		Ip:      connect.Ip,
//		KuId:    connect.KuId,
//		SteamId: connect.SteamId,
//	}
//	//获取最近的一条spwan记录和newComing
//	//playerLog := model.PlayerLog{Name: name, Role: spawn.Role, Action: action, ActionDesc: "", Time: time}
//	model.DB.Create(&playerLog)
//}
//
//func CollectSpawnRequestLog(text string) {
//	// Spawn request: wurt from 猜猜我是谁
//	arr := strings.Split(text, " ")
//	temp := strings.Replace(arr[0], " ", "", -1)
//	time := temp[:len(temp)-1]
//	role := strings.Replace(arr[3], " ", "", -1)
//	name := strings.Replace(arr[5], "\n", "", -1)
//
//	spawn := model.Spawn{Name: name, Role: role, Time: time}
//	model.DB.Create(&spawn)
//
//}
//
//func getSpawnRole(name string) *model.Spawn {
//	spawn := new(model.Spawn)
//	model.DB.Where("name LIKE ?", "%"+name+"%").Last(spawn)
//	return spawn
//}
//
//func getConnectInfo(name string) *model.Connect {
//	connect := new(model.Connect)
//	model.DB.Where("name LIKE ?", "%"+name+"%").Last(connect)
//	return connect
//}
//
//type Charset string
//
//const (
//	UTF8    = Charset("UTF-8")
//	GB18030 = Charset("GB18030")
//)
//
//func ConvertByte2String(byte []byte, charset Charset) string {
//
//	var str string
//	switch charset {
//	case GB18030:
//		decodeBytes, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
//		str = string(decodeBytes)
//	case UTF8:
//		fallthrough
//	default:
//		str = string(byte)
//	}
//	return str
//}
//
//func Tailf_server_chat_log(path string, level string) {
//	//fileName := "C:\\Users\\xm\\Documents\\Klei\\DoNotStarveTogether\\900587905\\Cluster_2\\Master\\server_chat_log.txt"
//	fileName := filepath.Join(path, level, "server_chat_log.txt")
//	log.Println("开始采集 server_chat_log, path:", fileName)
//	config := tail.Config{
//		ReOpen:    true,                                 // 重新打开
//		Follow:    true,                                 // 是否跟随
//		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件的哪个地方开始读
//		MustExist: false,                                // 文件不存在不报错
//		Poll:      true,
//	}
//	tails, err := tail.TailFile(fileName, config)
//	if err != nil {
//		log.Println("文件监听失败")
//	}
//	var (
//		line *tail.Line
//		ok   bool
//	)
//	for {
//		line, ok = <-tails.Lines
//		if !ok {
//			log.Println("文件监听失败")
//		}
//		//log.Println(line.Text)
//		CollectChatLog(line.Text)
//	}
//}
//
//func Tailf_server_log(path string, level string) {
//	//fileName := "C:\\Users\\xm\\Documents\\Klei\\DoNotStarveTogether\\900587905\\Cluster_2\\Master\\server_log.txt"
//	fileName := filepath.Join(path, level, "server_log.txt")
//	log.Println("开始采集 server_log, path:", fileName)
//	config := tail.Config{
//		ReOpen:    true,                                 // 重新打开
//		Follow:    true,                                 // 是否跟随
//		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件的哪个地方开始读
//		MustExist: false,                                // 文件不存在不报错
//		Poll:      true,
//	}
//	tails, err := tail.TailFile(fileName, config)
//	if err != nil {
//		log.Println("文件监听失败")
//	}
//	var (
//		line *tail.Line
//		ok   bool
//	)
//	var perLine string = ""
//	var start string = ""
//	first := true
//	connection := false
//	i := 0
//	var connect model.Connect
//	for {
//		line, ok = <-tails.Lines
//		if !ok {
//			log.Println("文件监听失败")
//		}
//
//		text := line.Text
//		perLine = text
//
//		if first {
//			start = text
//			//fmt.Println("日志", text)
//			first = false
//		}
//
//		//解析 时间
//		if find := strings.Contains(text, "# Generating"); find {
//			fmt.Println("房间结束了", start, perLine)
//		}
//		if find := strings.Contains(text, "Spawn request"); find {
//			CollectSpawnRequestLog(text)
//		}
//
//		//New incoming connection
//		if find := strings.Contains(text, "New incoming connection"); find {
//			connection = true
//			connect = model.Connect{}
//		}
//		if connection {
//			if i > 5 {
//				i = 0
//				connection = false
//			} else {
//				//do
//
//				if i == 1 {
//					// 解析 ip
//					fmt.Println(1, text)
//					str := strings.Split(text, " ")
//					if len(str) < 5 {
//						log.Println("[EROOR] str 解析错误: ", str)
//					} else {
//						var ip string
//						if strings.Contains(text, "[LAN]") {
//							ip = str[5]
//						} else {
//							ip = str[4]
//						}
//						connect.Ip = ip
//						fmt.Println("ip", ip)
//					}
//
//				} else if i == 3 {
//					fmt.Println(3, text)
//					// 解析 KU 和 用户名
//					str := strings.Split(text, " ")
//					if len(str) <= 4 {
//						log.Println("[EROOR] str 解析错误: ", str)
//					} else {
//						ku := str[3]
//						ku = ku[1 : len(ku)-1]
//						name := str[4]
//						connect.Name = name
//						connect.KuId = ku
//						fmt.Println("ku", ku, "name", name)
//					}
//				} else if i == 4 {
//					fmt.Println(4, text)
//					// 解析 steamId
//					str := strings.Split(text, " ")
//					if len(str) < 4 {
//						log.Println("[EROOR] str 解析错误: ", str)
//					} else {
//						steamId := str[4]
//						steamId = steamId[1 : len(steamId)-1]
//						fmt.Println("steamId", steamId)
//
//						//记录
//						connect.SteamId = steamId
//						model.DB.Create(&connect)
//					}
//				}
//				i++
//			}
//		}
//
//	}
//}
//
//func Tailf_server_log2(path string, level string) {
//	fileName := filepath.Join(path, level, "server_log.txt")
//	log.Println("开始采集 server_log, path:", fileName)
//	config := tail.Config{
//		ReOpen:    true,                                 // 重新打开
//		Follow:    true,                                 // 是否跟随
//		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件的哪个地方开始读
//		MustExist: false,                                // 文件不存在不报错
//		Poll:      true,
//	}
//	tails, err := tail.TailFile(fileName, config)
//	if err != nil {
//		log.Println("文件监听失败", err)
//	}
//	var (
//		which        = 0
//		isNewConnect = false
//		connect      model.Connect
//	)
//	for {
//		line, ok := <-tails.Lines
//		if !ok {
//			log.Println("文件读取失败", err)
//			time.Sleep(time.Second)
//		} else {
//			parseLog(line, &which, &isNewConnect, &connect)
//		}
//	}
//}
//
//func parseLog(line *tail.Line, which *int, isNewConnect *bool, connect *model.Connect) {
//
//	defer func() {
//		if err := recover(); err != nil {
//			log.Println("玩家日志解析异常:", err)
//		}
//	}()
//
//	text := line.Text
//	if find := strings.Contains(text, "Spawn request"); find {
//		CollectSpawnRequestLog(text)
//	}
//	//New incoming connection
//	if find := strings.Contains(text, "New incoming connection"); find {
//		*isNewConnect = true
//		connect = &model.Connect{}
//		*which = 0
//	}
//	if *isNewConnect {
//		if *which == 1 {
//			// 解析 ip
//			fmt.Println(1, text)
//			str := strings.Split(text, " ")
//			if len(str) < 5 {
//				log.Println("[EROOR] str 解析错误: ", str)
//				connect.Ip = ""
//			} else {
//				var ip string
//				if strings.Contains(text, "[LAN]") {
//					ip = str[5]
//				} else {
//					ip = str[4]
//				}
//				connect.Ip = ip
//				fmt.Println("ip", ip)
//			}
//
//		} else if *which == 3 {
//			fmt.Println(3, text)
//			// 解析 KU 和 用户名
//			str := strings.Split(text, " ")
//			if len(str) <= 4 {
//				log.Println("[EROOR] str 解析错误: ", str)
//			} else {
//				ku := str[3]
//				ku = ku[1 : len(ku)-1]
//				name := str[4]
//				connect.Name = name
//				connect.KuId = ku
//				fmt.Println("ku", ku, "name", name)
//			}
//		} else if *which == 4 {
//			fmt.Println(4, text)
//			// 解析 steamId
//			str := strings.Split(text, " ")
//			if len(str) < 4 {
//				log.Println("[EROOR] str 解析错误: ", str)
//			} else {
//				steamId := str[4]
//				steamId = steamId[1 : len(steamId)-1]
//				fmt.Println("steamId", steamId)
//
//				//记录
//				connect.SteamId = steamId
//				model.DB.Create(&connect)
//			}
//		}
//
//		*which = *which + 1
//	}
//}
