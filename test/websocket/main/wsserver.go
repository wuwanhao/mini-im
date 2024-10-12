package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// 创建一个 upGrader
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	http.HandleFunc("/hello", wsUpGrader)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Println("start serve err %v")
	}
}

func wsUpGrader(w http.ResponseWriter, r *http.Request) {
	// 将 http 连接升级为 main 连接，conn 代表返回的 main 连接对象
	conn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// 释放连接
	defer conn.Close()

	// 阻塞，等待消息处理
	for {
		// 接受消息
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("服务端收到消息，messageType:", messageType, ", message:", string(message))

		// 发送消息
		err = conn.WriteMessage(messageType, []byte("pong"))
		if err != nil {
			log.Println(err)
			return
		}
	}
}
