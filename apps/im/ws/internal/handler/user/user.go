package user

import (
	"app/apps/im/ws/internal/svc"
	websocketx "app/apps/im/ws/websocket"
	"github.com/gorilla/websocket"
)

func OnLine(svc *svc.ServiceContext) websocketx.HandleFunc {
	return func(srv *websocketx.Server, conn *websocket.Conn, msg *websocketx.Message) {
		uids := srv.GetUsers()
		u := srv.GetUsers(conn)
		// 给这条连接中所有的用户发送消息
		err := srv.Send(websocketx.NewMessage(u[0]+"上线了！", uids), conn)
		srv.Info(err)
	}
}
