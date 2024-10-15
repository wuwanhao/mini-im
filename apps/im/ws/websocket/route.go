package websocket

import "github.com/gorilla/websocket"

// 路由定义
type Route struct {
	Method  string     // 路由路径
	handler HandleFunc // 对应的处理器
}

type HandleFunc func(srv *Server, conn *websocket.Conn, msg *Message)
