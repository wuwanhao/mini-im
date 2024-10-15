package websocket

import (
	"fmt"
	"net/http"
	"time"
)

type Authentication interface {
	Auth(s *Server, w http.ResponseWriter, r *http.Request) bool // 认证的方法
	UserId(r *http.Request) string                               // 从 http 请求中获取 userId
}

// 实现接口
type authentication struct {
}

func (a *authentication) Auth(s *Server, w http.ResponseWriter, r *http.Request) bool {
	return true
}

func (a *authentication) UserId(r *http.Request) string {
	query := r.URL.Query()
	if query != nil && query["userId"] != nil {
		return fmt.Sprintf("%v", query["userId"])
	}
	return fmt.Sprintf("%v", time.Now().UnixMilli())
}
