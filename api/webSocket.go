package api

import (
	"dst-admin-go/utils/fileUtils"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hpcloud/tail"
)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Poll file for changes with this period.
	filePeriod = 10 * time.Second
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func HandlerWS(ctx *gin.Context) {
	w := ctx.Writer
	r := ctx.Request
	wsTail(w, r)
}

func wsTail(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Http uprader websocket error")
		log.Println(err.Error())
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}
	var done int32 = 0
	go heartbeat(ws)
	wsReader(ws, done)
}

func heartbeat(ws *websocket.Conn) {
	pingTicker := time.NewTicker(pingPeriod)
	fileTicker := time.NewTicker(filePeriod)

	defer func() {
		pingTicker.Stop()
		fileTicker.Stop()
		ws.Close()
		log.Println("ws 心跳检测退出")
	}()

	for {
		//从定时器中获取数据
		<-pingTicker.C
		log.Println("pingTicker")
		ws.SetWriteDeadline(time.Now().Add(writeWait))
		if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
			log.Println(err)
			return
		}
	}

	// for {
	// 	select {
	// 	case <-pingTicker.C:
	// 		log.Println("pingTicker")
	// 		ws.SetWriteDeadline(time.Now().Add(writeWait))
	// 		if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
	// 			log.Println(err)
	// 			return
	// 		}
	// 	}
	// }

}

func wsReader(ws *websocket.Conn, done int32) {

	defer ws.Close()
	log.Println("ws close")
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		log.Printf("type: %d", messageType)
		log.Printf("message: %s", message)
		msg := string(message)
		//监听文件
		if messageType == 1 && strings.HasPrefix(msg, "tailf") {
			split := strings.Split(msg, " ")
			if len(split) != 2 {
				continue
			}
			filename := strings.TrimSpace(split[1])
			log.Printf("tailf " + filename)

			//返回前100条数据
			content, err := fileUtils.ReverseRead(filename, 100)
			if err == nil {
				for _, line := range content {
					if err := ws.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
						log.Println(err)
						return
					}
				}
			}
			//tailf 文件
			go tailf(ws, filename, done)
		}

		if messageType == 1 && msg == "byte" {
			atomic.AddInt32(&done, 1)
			break
		}
	}
}

func tailf(ws *websocket.Conn, filename string, done int32) {
	fileName := filename
	config := tail.Config{
		ReOpen:    true,                                 // 重新打开
		Follow:    true,                                 // 是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件的哪个地方开始读
		MustExist: false,                                // 文件不存在不报错
		Poll:      true,
	}
	tails, err := tail.TailFile(fileName, config)
	if err != nil {
		fmt.Println("tail file failed, err:", err)
		return
	}
	var (
		line *tail.Line
		ok   bool
	)
	for {
		line, ok = <-tails.Lines
		if !ok {
			fmt.Printf("tail file close reopen, filename:%s\n", tails.Filename)
			time.Sleep(time.Second)
			continue
		}
		ws.SetWriteDeadline(time.Now().Add(writeWait))
		if err := ws.WriteMessage(websocket.TextMessage, []byte(line.Text)); err != nil {
			log.Println(err)
			return
		}
		// fmt.Println("line:", line.Text)
		if atomic.CompareAndSwapInt32(&done, 1, 1) {
			break
		}

	}
	log.Println("tailf close")
}
