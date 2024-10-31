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

	// 用户-连接关系映射
	connToUser map[*websocket.Conn]string
	userToConn map[string]*websocket.Conn

	patten string
	opt    option
}

func NewServer(addr string, opts ...Options) *Server {
	// 这里使用了函数选项设计模式，进行服务的初始化工作
	opt := newOption(opts...)
	return &Server{
		routes:         make(map[string]HandleFunc),
		addr:           addr,
		upGrader:       websocket.Upgrader{
			// 设置 websocket 允许跨域
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		authentication: opt.Authentication,
		Logger:         logx.WithContext(context.Background()),
		connToUser:     make(map[*websocket.Conn]string),
		userToConn:     make(map[string]*websocket.Conn),
		patten:         opt.pattern,
		opt:            opt,
	}
}

// Server handle 将 http 升级为 websocket
func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {
	// 注册异常处理函数
	defer func() {
		if r := recover(); r != nil {
			s.Errorf("Server Handler recover err: %v", r)
		}
	}()

	// 升级连接，由http升级为 websocket
	conn, err := s.upGrader.Upgrade(w, r, nil)
	if err != nil {
		s.Errorf("upgrade http conn err: %v", err)
		return
	}

	// 身份认证
	if !s.opt.Authentication.Auth(w, r) {
		s.Info("auth failed")
		conn.WriteMessage(websocket.TextMessage, []byte("认证失败"))
		conn.Close()
		return
	}

	// 添加连接记录，会有并发问题
	s.AddConn(conn, r)
	// 用一个协程来处理连接
	go s.HandleConn(conn)
}

// 处理连接
func (s *Server) HandleConn(conn *websocket.Conn) {
	// 阻塞处理
	for {
		// 获取请求消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("websocket read message err: %v", err)
			// 获取消息失败的时候关闭连接
			s.Close(conn)
			return
		}

		// 拿到 websocket 请求参数
		message := &Message{}
		if err = json.Unmarshal(msg, message); err != nil {
			s.Errorf("json unmarshal message err: %v, msg: %v", err, msg)
			s.Close(conn)
			return
		}

		// 根据请求的 method 分发路由并执行
		if handler, ok := s.routes[message.Method]; ok {
			handler(s, conn, message)
		} else {
			err := conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("不存在的请求方法：%v 请检查", message.Method)))
			if err != nil {
				return
			}
		}
	}
}

// 启动服务器
func (s *Server) Start() {
	http.HandleFunc(s.patten, s.ServerWs)
	s.Info(http.ListenAndServe(s.addr, nil))
}

// 停止服务器
func (s *Server) Stop() {
	fmt.Println("STOP SERVER")
}

// 关闭某个 websocket 连接
func (s *Server) Close(conn *websocket.Conn) {
	conn.Close()

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	// 删除对应的 uid-conn 连接
	uid := s.connToUser[conn]
	delete(s.connToUser, conn)
	delete(s.userToConn, uid)
}

// 添加路由
func (s *Server) AddRoutes(routes []Route) {
	for _, route := range routes {
		s.routes[route.Method] = route.Handler
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
	// 这里使用读锁
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

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
	// 这里使用读锁
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	return s.userToConn[uid]
}

// 根据用户 ID 批量获取连接
func (s *Server) GetConns(uids ...string) []*websocket.Conn {
	if len(uids) == 0 {
		return nil
	}
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	res := make([]*websocket.Conn, 0, len(uids))
	for _, uid := range uids {
		res = append(res, s.userToConn[uid])
	}
	return res
}

// 通过 userId 发送消息
func (s *Server) SendByUserId(msg interface{}, userIds ...string) error {
	if len(userIds) == 0 {
		return nil
	}

	return s.Send(msg, s.GetConns(userIds...)...)
}

// 发送消息
func (s *Server) Send(msg interface{}, conn ...*websocket.Conn) error {
	if len(conn) == 0 {
		return nil
	}

	// 序列化消息
	data, err := json.Marshal(msg)
	if err != nil {
		return nil
	}

	// 将序列化后的消息信息发送到所有的连接中
	for _, c := range conn {
		if err := c.WriteMessage(websocket.TextMessage, data); err != nil {
			return err
		}
	}

	return nil

}
