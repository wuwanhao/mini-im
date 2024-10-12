package main

import (
	"github.com/gorilla/websocket"
	"log"
)

func main() {
	url := "ws://localhost:8080/hello"

	// 使用默认拨号器，向服务端发起连接
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// 关闭连接
	defer conn.Close()

	// 发送消息
	go func() {
		for {
			err := conn.WriteMessage(websocket.TextMessage, []byte("ping"))
			if err != nil {
				log.Println(err)
				return
			}
		}
	}()

	// 接受消息
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("服务端收到消息，message:", string(message))
	}

}
