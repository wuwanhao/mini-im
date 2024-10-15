package handler

import (
	"app/apps/im/ws/internal/handler/user"
	"app/apps/im/ws/internal/svc"
	"app/apps/im/ws/websocket"
)

func RegisterHandlers(srv *websocket.Server, svc *svc.ServiceContext) {
	srv.AddRoutes([]websocket.Route{
		{
			Method:  "user.OnLine",
			Handler: user.OnLine(svc),
		},
	})
}
