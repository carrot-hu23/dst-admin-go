package chatgpt

import (
	"dst-admin-go/entity"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	HTTPS_PROXY    = "http://127.0.0.1:7890"
	OPENAI_API_URL = "https://api.openai.com/v1/chat/completions"
	Model          = "gpt-3.5-turbo"
)

type ModifiedResponse struct {
	reader io.Reader
}

type MessagesBody struct {
	Messages []ChatMessage `json:"messages"`
}

func (m *ModifiedResponse) Read(p []byte) (int, error) {
	fmt.Println("p[]byte", string(p))
	n, err := m.reader.Read(p)
	if err == nil {
		// 在读取数据时对数据进行修改
		upper := strings.ToUpper(string(p[:n]))
		copy(p[:n], []byte(upper))
	}
	return n, err
}

func ChatGpt(text string, f func(message string)) {
	chatMessages := []ChatMessage{
		{Role: "system", Content: entity.Config.Prompt},
		{Role: "user", Content: text},
	}
	response := Post(chatMessages)
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
		return
	}
	content := data["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	f(content)
}

func Post(messages []ChatMessage) *http.Response {
	http_proxy, _ := url.Parse(HTTPS_PROXY)

	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(http_proxy),
		},
	}

	payload := strings.NewReader(generatePayload(messages))

	req, _ := http.NewRequest("POST", OPENAI_API_URL, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+entity.Config.OPENAI_API_KEY)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	return resp
}

func generatePayload(message []ChatMessage) string {

	body := map[string]interface{}{
		"model":       Model,
		"messages":    message,
		"temperature": 0.6,
		// "stream":      true,
	}
	data, err := json.Marshal(body)
	if err != nil {
		fmt.Println("josn parse error: ", err)
	}
	return string(data)
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
