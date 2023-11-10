package kook

import (
	"github.com/gin-gonic/gin"
	"log"
)

type KookBotApi struct{}

type Body struct {
	S int `json:"s"`
	D D   `json:"d"`
}
type D struct {
	Type        int    `json:"type"`
	ChannelType string `json:"channel_type"`
	Challenge   string `json:"challenge"`
	VerifyToken string `json:"verify_token"`
}

type CardMessage struct {
	Type    string               `json:"type"`
	Theme   string               `json:"theme"`
	Size    string               `json:"size"`
	Modules []CardMessageModules `json:"modules"`
}

type CardMessageModules struct {
	Type     string `json:"type"`
	Text     Text   `json:"text"`
	Elements []struct {
		Type  string `json:"type"`
		Src   string `json:"src"`
		Theme string `json:"theme"`
		Value string `json:"value"`
		Text  struct {
			Type    string `json:"type"`
			Content string `json:"content"`
		} `json:"text"`
	} `json:"elements"`
	Title string `json:"title"`
	Src   string `json:"src"`
	Size  string `json:"size"`
}

type Text struct {
	Type    string  `json:"type"`
	Content string  `json:"content"`
	Cols    int     `json:"cols"`
	Fields  []Field `json:"fields"`
}

type Field struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type Elements struct {
	Type  string `json:"type"`
	Src   string `json:"src"`
	Theme string `json:"theme"`
	Value string `json:"value"`
	Text  Text   `json:"text"`
}

func (k *KookBotApi) AuthKookWebHook(ctx *gin.Context) {

	var requestData map[string]interface{}
	if err := ctx.Bind(&requestData); err != nil {
		log.Println(err)
	} else {
		log.Println(requestData)
	}

	var requestBody Body

	if err := ctx.BindJSON(&requestBody); err != nil {
		ctx.JSON(400, gin.H{"error": "无效的请求"})
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

	var cardMessages []CardMessage
	cardMessages = append(cardMessages, CardMessage{
		Type:  "card",
		Theme: "secondary",
		Size:  "lg",
		Modules: []CardMessageModules{
			{
				Type: "header",
				Text: Text{
					Type:    "plain-text",
					Content: "指令导航",
				},
			},
			{
				Type: "divider",
			},
			{
				Type: "section",
				Text: Text{
					Type: "paragraph",
					Cols: 3,
					Fields: []Field{
						{
							Type:    "kmarkdown",
							Content: "**指令**\n/query\n/info\n/player",
						},
						{
							Type:    "kmarkdown",
							Content: "**操作**\n/query\n/info\n/player",
						},
						{
							Type:    "kmarkdown",
							Content: "**作用**\n/query\n/info\n/player",
						},
					},
				},
			},
		},
	})

	ctx.JSON(200, cardMessages)
}
