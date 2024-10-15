package websocket

// ws 消息结构定义
type Message struct {
	Method string      `json:"method,omitempty"`
	UserId string      `json:"user_id,omitempty"`
	FromId string      `json:"from_id,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}
