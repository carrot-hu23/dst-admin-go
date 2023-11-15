package kook

import (
	"dst-admin-go/bot/kook/http"
	"dst-admin-go/bot/kook/message"
	"dst-admin-go/service"
	"dst-admin-go/utils/dstConfigUtils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"strings"
	"time"
)

var gameArchive = service.GameArchive{}
var playerService = service.PlayerService{}

type KookBotApi struct{}

type D struct {
	Type        int    `json:"type"`
	ChannelType string `json:"channel_type"`
	Challenge   string `json:"challenge"`
	VerifyToken string `json:"verify_token"`

	TargetId     string                 `json:"target_id"`
	Nonce        string                 `json:"nonce"`
	MsgId        string                 `json:"msg_id"`
	MsgTimestamp int64                  `json:"msg_timestamp"`
	FromType     int                    `json:"from_type"`
	Content      string                 `json:"content"`
	Extra        map[string]interface{} `json:"extra"`
}

type RequestData struct {
	S  int `json:"s"`
	Sn int `json:"sn"`
	D  D   `json:"d"`
}

func (k *KookBotApi) AuthKookWebHook(ctx *gin.Context) {

	var requestBody RequestData
	if err := ctx.ShouldBind(&requestBody); err != nil {
		log.Println(err)
		ctx.JSON(401, gin.H{"error": "无效的请求"})
		return
	}
	// 表示这是一个验证请求
	if requestBody.D.ChannelType == "WEBHOOK_CHALLENGE" {
		challenge := requestBody.D.Challenge
		responseData := gin.H{
			"challenge": challenge,
		}
		ctx.JSON(200, responseData)
		return
	}
	log.Println("requestBody: ", requestBody)

	if requestBody.D.ChannelType != "GROUP" {
		// 无效请求 记录下
		ctx.JSON(200, gin.H{"error": "无效的请求"})
		return
	}
	handleReq(&requestBody, ctx)
	ctx.JSON(200, "")

}

func getCommand(requestData *RequestData) string {
	split := strings.Split(requestData.D.Content, "/")
	if len(split) != 2 {
		return ""
	}
	return split[1]
}

const BaseUrl = "https://www.kookapp.cn"
const Token = "1/MjQ5MzU=/xigN7Aa0FDTkO64s3ksQkw=="

var messageClient = http.NewMessageClient(BaseUrl, Token)
var commandMap = map[string]func(data *RequestData){}

func init() {

	commandMap["帮助"] = HandleHelp

	commandMap["ping"] = HandlePing
	commandMap["房间信息"] = HandleHomeInfo
	commandMap["玩家列表"] = HandlePlayerList

}

func handleReq(requestData *RequestData, ctx *gin.Context) {

	command := requestData.D.Content
	if f, ok := commandMap[command]; ok {
		f(requestData)
	}

}

func HandleHelp(requestData *RequestData) {
	// 返回导航

	//var cardMessages []message.CardMessage
	//cardMessages = append(cardMessages, message.CardMessage{
	//	Type:  "card",
	//	Theme: "secondary",
	//	Size:  "lg",
	//	Modules: []message.CardMessageModules{
	//		{
	//			Type: "header",
	//			Text: message.Text{
	//				Type:    "plain-text",
	//				Content: "指令导航",
	//			},
	//		},
	//		{
	//			Type: "divider",
	//		},
	//		{
	//			Type: "section",
	//			Text: message.Text{
	//				Type: "paragraph",
	//				Cols: 3,
	//				Fields: []message.Field{
	//					{
	//						Type:    "kmarkdown",
	//						Content: "**指令**\n/query\n/info\n/player",
	//					},
	//					{
	//						Type:    "kmarkdown",
	//						Content: "**操作**\n/query\n/info\n/player",
	//					},
	//					{
	//						Type:    "kmarkdown",
	//						Content: "**作用**\n/query\n/info\n/player",
	//					},
	//				},
	//			},
	//		},
	//	},
	//})
	//
	//content, err := json.Marshal(cardMessages)
	//if err != nil {
	//	log.Println(err)
	//}

	_, err := messageClient.Create(message.Kmarkdown_Type, requestData.D.TargetId, "**饥荒机器人使用帮助**\n---\n          \n        1️⃣ 输入 `帮助` 获取指令帮助 \n        2️⃣ 输入 `房间信息` 获取房间信息\n        3️⃣ 输入 `玩家列表` 获取房间玩家列表\n        4️⃣ 输入 `查询房间+房间名称+页数` 查询房间\n        5️⃣ 输入 `ping` 返回pong \n        6️⃣ 禁止带节奏的言论。\n        7️⃣ 对管理员处理方式有异议请私信其他管理员解决。\n        8️⃣ 被人帮助，请说谢谢，组队一起语音避免沟通不当。\n          \n ", requestData.D.MsgId, "", "")
	if err != nil {
		log.Println(err)
	}
}

func HandlePing(requestData *RequestData) {
	_, err := messageClient.Create(message.Text_Type, requestData.D.TargetId, "pong date: "+time.Now().Format("2006-01-02 15:04:05"), requestData.D.MsgId, "", "")
	if err != nil {
		log.Println(err)
	}
}

func HandleHomeInfo(requestData *RequestData) {
	config := dstConfigUtils.GetDstConfig()
	clusterName := config.Cluster
	start := time.Now() // 获取当前时间
	archive := gameArchive.GetGameArchive(clusterName)
	elapsed := time.Since(start)
	log.Println("该函数执行完成耗时：", elapsed)
	var cardMessages []message.CardMessage
	cardMessages = append(cardMessages, message.CardMessage{
		Type:  "card",
		Theme: "secondary",
		Size:  "lg",
		Color: "#516cd9",
		Modules: []message.CardMessageModules{
			{
				Type: "header",
				Text: message.Text{
					Type:    "plain-text",
					Content: "房间信息",
				},
			},
			{
				Type: "divider",
			},
			{
				Type: "section",
				Text: message.Text{
					Type: "paragraph",
					Cols: 3,
					Fields: []message.Field{
						{
							Type:    "kmarkdown",
							Content: "**房间名称**\n" + archive.ClusterName,
						},
						{
							Type:    "kmarkdown",
							Content: "**人数**\n" + strconv.Itoa(len(archive.Players)) + "/" + strconv.Itoa(archive.MaxPlayers),
						},
						{
							Type:    "kmarkdown",
							Content: "**天数**\n" + strconv.Itoa(archive.Meta.Clock.Cycles) + " 天",
						},
						{
							Type:    "kmarkdown",
							Content: "**季节**\n" + archive.Meta.Seasons.Season + "(" + strconv.Itoa(archive.Meta.Seasons.ElapsedDaysInSeason) + "/" + strconv.Itoa(archive.Meta.Seasons.ElapsedDaysInSeason+archive.Meta.Seasons.RemainingDaysInSeason) + ")",
						},
						{
							Type:    "kmarkdown",
							Content: "**模式**\n" + archive.GameMod,
						},
						{
							Type:    "kmarkdown",
							Content: "**模组**\n" + strconv.Itoa(archive.Mods),
						},
						{
							Type:    "kmarkdown",
							Content: "**版本**\n" + strconv.FormatInt(archive.Version, 10) + "/" + strconv.FormatInt(archive.LastVersion, 10),
						},
					},
				},
			},
		},
	})

	content, err := json.Marshal(cardMessages)
	if err != nil {
		log.Println(err)
	}
	data, err := messageClient.Create(message.Card_Type, requestData.D.TargetId, string(content), requestData.D.MsgId, "", "")
	if err != nil {
		log.Println(err)
	}

	// 进行类型断言将接口对象转换为map类型
	if m, ok := data.(map[string]interface{}); ok {
		// 访问code和message值
		code := m["code"].(float64)
		msg := m["message"].(string)
		log.Println(msg)
		if code != 200 {
			text := "**# 房间信息**\n---\n房间名称: " + archive.ClusterName + "\n人数: " + strconv.Itoa(len(archive.Players)) + "/" + strconv.Itoa(archive.MaxPlayers) + "\n模组: " + strconv.Itoa(archive.Mods) + "\n天数: " + strconv.Itoa(archive.Meta.Clock.Cycles) + " 天" + "\n季节: " + archive.Meta.Seasons.Season + "(" + strconv.Itoa(archive.Meta.Seasons.ElapsedDaysInSeason) + "/" + strconv.Itoa(archive.Meta.Seasons.ElapsedDaysInSeason+archive.Meta.Seasons.RemainingDaysInSeason) + ")" + "\n模式: " + archive.GameMod + "\n版本: " + strconv.FormatInt(archive.Version, 10) + "/" + strconv.FormatInt(archive.LastVersion, 10)
			// messageClient.Create(message.Text_Type, requestData.D.TargetId, "http code: "+strconv.Itoa(int(code))+" msg: "+msg, requestData.D.MsgId, "", "")
			messageClient.Create(message.Kmarkdown_Type, requestData.D.TargetId, text, requestData.D.MsgId, "", "")

		}
	} else {
		fmt.Println("Invalid type for data")
	}

}

func HandlePlayerList(requestData *RequestData) {
	config := dstConfigUtils.GetDstConfig()
	clusterName := config.Cluster
	playerlist := playerService.GetPlayerList(clusterName, "Master")

	var nameset string
	var roleset string
	var dayset string
	for i := range playerlist {
		nameset = nameset + playerlist[i].Name + "(" + playerlist[i].KuId + ")" + "\n"
		roleset = roleset + playerlist[i].Role + "\n"
		dayset = dayset + playerlist[i].Day + " 天" + "\n"

	}

	var cardMessages []message.CardMessage
	cardMessages = append(cardMessages, message.CardMessage{
		Type:  "card",
		Theme: "secondary",
		Size:  "lg",
		Color: "#516cd9",
		Modules: []message.CardMessageModules{
			{
				Type: "header",
				Text: message.Text{
					Type:    "plain-text",
					Content: "玩家列表",
				},
			},
			{
				Type: "divider",
			},
			{
				Type: "section",
				Text: message.Text{
					Type: "paragraph",
					Cols: 3,
					Fields: []message.Field{
						{
							Type:    "kmarkdown",
							Content: "**玩家**\n" + nameset,
						},
						{
							Type:    "kmarkdown",
							Content: "**角色**\n" + roleset,
						},
						{
							Type:    "kmarkdown",
							Content: "**天数**\n" + dayset,
						},
					},
				},
			},
		},
	})

	content, err := json.Marshal(cardMessages)
	if err != nil {
		log.Println(err)
	}

	data, err := messageClient.Create(message.Card_Type, requestData.D.TargetId, string(content), requestData.D.MsgId, "", "")
	if err != nil {
		log.Println(err)
	}

	// 进行类型断言将接口对象转换为map类型
	if m, ok := data.(map[string]interface{}); ok {
		log.Println("m", m)
		// 访问code和message值
		code := m["code"].(float64)
		msg := m["message"].(string)
		if code != 200 {
			messageClient.Create(message.Text_Type, requestData.D.TargetId, "http code: "+strconv.Itoa(int(code))+" msg: "+msg, requestData.D.MsgId, "", "")
		}
	} else {
		fmt.Println("Invalid type for data")
	}

}
