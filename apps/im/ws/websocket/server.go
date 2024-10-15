package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"sync"
)

// websocket 服务对象
type Server struct {
	authentication Authentication
	sync.RWMutex
	routes   map[string]HandleFunc // 路由地址以及处理方法
	addr     string                // 服务地址
	upGrader websocket.Upgrader    // websocket升级器
	logx.Logger

	// 用户连接信息
	connToUser map[*websocket.Conn]string
	userToConn map[string]*websocket.Conn

	opt option
}

func NewServer(addr string, opts ...Options) *Server {
	// 这里使用了函数选项设计模式
	opt := newOption(opts...)
	return &Server{
		routes:     make(map[string]HandleFunc),
		addr:       addr,
		upGrader:   websocket.Upgrader{},
		Logger:     logx.WithContext(context.Background()),
		connToUser: make(map[*websocket.Conn]string),
		userToConn: make(map[string]*websocket.Conn),
		opt:        opt,
	}
}
	
// Server handle 将 http 升级为 websocket
func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {
	// 注册异常处理函数
	defer func() {
		if r := recover(); r != nil {
			s.Errorf("Server handler recover err: %v", r)
		}
	}()

	// 身份认证
	if !s.opt.Authentication.Auth(s, w, r) {
		s.Info("auth failed")
		return
	}
	// 升级连接，由http升级为 websocket
	conn, err := s.upGrader.Upgrade(w, r, nil)
	if err != nil {
		s.Errorf("upgrade http conn err: %v", err)
	}

	// 添加连接记录，会有并发问题
	s.AddConn(conn, r)
	// 处理连接
	go s.HandleConn(conn)
}

// 处理连接
func (s *Server) HandleConn(conn *websocket.Conn) {
	// 阻塞处理
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("websocket read message err: %v", err)
			// 关闭连接
			s.Close(conn)
			return
		}

		// 拿到 websocket request param
		message := &Message{}
		json.Unmarshal(msg, message)

		// 判断是否有对应的处理器来处理此请求
		if handler, ok := s.routes[message.Method]; ok {
			handler(s, conn, message)
		} else {
			conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("不存在的请求方法：%v 请检查", message.Method)))
		}
	}
}

// 启动服务器
func (s *Server) Start() {
	http.HandleFunc("/ws", s.ServerWs)
	http.ListenAndServe(s.addr, nil)
}

// 停止服务器
func (s *Server) Stop() {
	fmt.Println("STOP SERVER")
}

// 关闭某个 websocket 连接
func (s *Server) Close(conn *websocket.Conn) {
	conn.Close()
}

// 添加路由
func (s *Server) AddRoutes(routes []Route) {
	for _, route := range routes {
		s.routes[route.Method] = route.handler
	}
}

// 添加连接
func (s *Server) AddConn(conn *websocket.Conn, r *http.Request) {
	uid := s.authentication.UserId(r)
	// 此处是 map 的写操作，会存在并发问题
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	s.connToUser[conn] = uid
	s.userToConn[uid] = conn
}

// 根据连接获取用户
func (s *Server) GetUsers(conns ...*websocket.Conn) []string {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	var res []string
	if len(conns) == 0 {
		// 获取全部
		res = make([]string, 0, len(s.connToUser))
		// conn => uid
		for _, uid := range s.connToUser {
			res = append(res, uid)
		}
	} else {
		// 获取给定连接的用户
		res = make([]string, 0, len(conns))
		for _, conn := range conns {
			res = append(res, s.connToUser[conn])
		}
	}

	return res
}

// 根据用户获取连接
func (s *Server) GetConn(uid string) *websocket.Conn {
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()
	return s.userToConn[uid]
}
