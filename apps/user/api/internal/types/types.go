// Code generated by goctl. DO NOT EDIT.
package types

type User struct {
	Id       string `json:"id"`
	NickName string `json:"nickname"`
	Phone    string `json:"phone"`
	Sex      byte   `json:"sex"`
	Avatar   string `json:"avatar"`
}

type RegisterReq struct {
	Phone    string `json:"phone"`
	NickName string `json:"nickname"`
	Password string `json:"password"`
	Sex      byte   `json:"sex"`
	Avatar   string `json:"avatar"`
}

type RegisterResp struct {
	Token  string `json:"token"`
	Expire int64  `json:"expire"`
}

type LoginReq struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type LoginResp struct {
	Token  string `json:"token"`
	Expire int64  `json:"expire"`
}

type UserInfoReq struct {
}

type UserInfoResp struct {
	Info User `json:"info"`
}
