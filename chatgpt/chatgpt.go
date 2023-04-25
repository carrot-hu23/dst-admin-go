package chatgpt

import (
	"container/list"
	"dst-admin-go/entity"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

const (
	HTTPS_PROXY    = "http://127.0.0.1:7890"
	OPENAI_API_URL = "https://api.openai.com/v1/chat/completions"
	Model          = "gpt-3.5-turbo"
)

var History = NewLRUCache(20)

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

func ChatGpt(user, text string, f func(message string)) {

	message := ChatMessage{Role: "user", Content: text}
	messages := History.AddMessage(user, message)

	// message := ChatMessage{Role: "user", Content: text}
	// messages := History.AddMessage(user, message)
	// chatMessages := []ChatMessage{
	// 	{Role: "system", Content: entity.Config.Prompt},
	// 	{Role: "user", Content: text},
	// }
	response := Post(messages)
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
	str := strings.ReplaceAll(content, "\n", "\\\\n")
	f(str)

	History.AddMessage(user, ChatMessage{Role: "system", Content: content})
}

func randApikey() string {
	s := entity.Config.OPENAI_API_KEY
	availableAPIKeys := strings.Split(s, ",")
	key := availableAPIKeys[rand.Intn(len(availableAPIKeys))]
	return key
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
	req.Header.Add("Authorization", "Bearer "+randApikey())

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

type LRUCache struct {
	capacity int
	cache    map[string]*list.Element
	list     *list.List
}

type Pair struct {
	Key   string
	Value []ChatMessage
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		list:     list.New(),
	}
}

func (c *LRUCache) Get(key string) ([]ChatMessage, bool) {
	if elem, ok := c.cache[key]; ok {
		c.list.MoveToFront(elem)
		return elem.Value.(*Pair).Value, true
	}
	return nil, false
}

func (c *LRUCache) Put(key string, value []ChatMessage) {
	if elem, ok := c.cache[key]; ok {
		// Move the element to the front of the list
		c.list.MoveToFront(elem)
		// Update the value of the element
		elem.Value.(*Pair).Value = value
	} else {
		// Add a new element to the front of the list
		elem := c.list.PushFront(&Pair{Key: key, Value: value})
		c.cache[key] = elem
		// Remove the least recently used element if the cache is full
		if c.list.Len() > c.capacity {
			tail := c.list.Back()
			delete(c.cache, tail.Value.(*Pair).Key)
			c.list.Remove(tail)
		}
	}
}

func (c *LRUCache) AddMessage(key string, message ChatMessage) []ChatMessage {
	if elem, ok := c.cache[key]; ok {
		// Move the element to the front of the list
		c.list.MoveToFront(elem)
		// Add the new message to the value of the element
		pair := elem.Value.(*Pair)
		pair.Value = append(pair.Value, message)
		// Limit the size of the slice to 20
		if len(pair.Value) > 20 {
			pair.Value = pair.Value[len(pair.Value)-20:]
		}
		return pair.Value
	} else {
		// Add a new element to the front of the list
		elem := c.list.PushFront(&Pair{Key: key, Value: []ChatMessage{message}})
		c.cache[key] = elem
		// Remove the least recently used element if the cache is full
		if c.list.Len() > c.capacity {
			tail := c.list.Back()
			delete(c.cache, tail.Value.(*Pair).Key)
			c.list.Remove(tail)
		}
		return []ChatMessage{message}
	}
}
