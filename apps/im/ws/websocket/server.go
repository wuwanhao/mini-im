package websocket

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

// websocket 服务对象
type Server struct {
	addr     string             // 服务地址
	upGrader websocket.Upgrader // websocket升级器
	logx.Logger
}

func NewServer(addr string) *Server {
	return &Server{
		addr:     addr,
		upGrader: websocket.Upgrader{},
		Logger:   logx.WithContext(context.Background()),
	}
}

func ServerWs(w http.ResponseWriter, r *http.Request) {}

func (s *Server) Start() {
	http.HandleFunc("/ws", ServerWs)
	http.ListenAndServe(s.addr, nil)
}

func (s *Server) Stop() {
	fmt.Println("STOP SERVER")
}
