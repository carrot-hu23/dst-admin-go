package schedule

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type ImageSchedule struct {
}

type ImageTask struct {
	Title string  `json:"title"`
	Url   string  `json:"url"`
	Date  float64 `json:"date"`
}

func (i *ImageSchedule) parseImageHtml(task ImageTask) {
	resp, err := http.Get(task.Url)
	if err != nil {
		log.Println(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	s := string(body)

	log.Println(s)
}

func (i *ImageSchedule) StartSchedule() []ImageTask {
	url := "https://steamcommunity-a.akamaihd.net/news/newsforapp/v0002/?appid=322330&count=10&maxlength=300&format=json"

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	s := string(body)
	var result map[string]interface{}
	err = json.Unmarshal([]byte(s), &result)

	newsitems := result["appnews"].(map[string]interface{})["newsitems"].([]interface{})
	var tasks []ImageTask
	for _, newsitem := range newsitems {
		if newsitem.(map[string]interface{})["feed_type"].(float64) == 1 {

			title := newsitem.(map[string]interface{})["title"].(string)
			url := newsitem.(map[string]interface{})["url"].(string)
			date := newsitem.(map[string]interface{})["date"].(float64)
			task := ImageTask{
				Title: title,
				Url:   url,
				Date:  date,
			}
			tasks = append(tasks, task)
		}
	}
	return tasks
}
