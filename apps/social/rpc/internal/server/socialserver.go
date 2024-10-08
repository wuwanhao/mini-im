// Code generated by goctl. DO NOT EDIT.
// Source: social.proto

package server

import (
	"context"

	"app/apps/social/rpc/internal/logic"
	"app/apps/social/rpc/internal/svc"
	"app/apps/social/rpc/rpc"
)

type SocialServer struct {
	svcCtx *svc.ServiceContext
	rpc.UnimplementedSocialServer
}

func NewSocialServer(svcCtx *svc.ServiceContext) *SocialServer {
	return &SocialServer{
		svcCtx: svcCtx,
	}
}

// 好友业务：请求添加好友、通过或拒绝申请、好有列表
func (s *SocialServer) FriendPutIn(ctx context.Context, in *rpc.FriendPutInReq) (*rpc.FriendPutInResp, error) {
	l := logic.NewFriendPutInLogic(ctx, s.svcCtx)
	return l.FriendPutIn(in)
}

func (s *SocialServer) FriendPutInHandle(ctx context.Context, in *rpc.FriendPutInHandleReq) (*rpc.FriendPutInHandleResp, error) {
	l := logic.NewFriendPutInHandleLogic(ctx, s.svcCtx)
	return l.FriendPutInHandle(in)
}

func (s *SocialServer) FriendList(ctx context.Context, in *rpc.FriendListReq) (*rpc.FriendListResp, error) {
	l := logic.NewFriendListLogic(ctx, s.svcCtx)
	return l.FriendList(in)
}

func (s *SocialServer) FriendPutInList(ctx context.Context, in *rpc.FriendPutInListReq) (*rpc.FriendPutInListResp, error) {
	l := logic.NewFriendPutInListLogic(ctx, s.svcCtx)
	return l.FriendPutInList(in)
}

// 群组业务：创建群组、修改群、群公告、申请群、加群请求列表、加群请求处理...
func (s *SocialServer) GroupCreate(ctx context.Context, in *rpc.GroupCreateReq) (*rpc.GroupCreateResp, error) {
	l := logic.NewGroupCreateLogic(ctx, s.svcCtx)
	return l.GroupCreate(in)
}

func (s *SocialServer) GroupPutIn(ctx context.Context, in *rpc.GroupPutInReq) (*rpc.GroupPutInResp, error) {
	l := logic.NewGroupPutInLogic(ctx, s.svcCtx)
	return l.GroupPutIn(in)
}

func (s *SocialServer) GroupPutInList(ctx context.Context, in *rpc.GroupPutInListReq) (*rpc.GroupPutInListResp, error) {
	l := logic.NewGroupPutInListLogic(ctx, s.svcCtx)
	return l.GroupPutInList(in)
}

func (s *SocialServer) GroupPutInHandle(ctx context.Context, in *rpc.GroupPutInHandleReq) (*rpc.GroupPutInHandleResp, error) {
	l := logic.NewGroupPutInHandleLogic(ctx, s.svcCtx)
	return l.GroupPutInHandle(in)
}

func (s *SocialServer) GroupList(ctx context.Context, in *rpc.GroupListReq) (*rpc.GroupListResp, error) {
	l := logic.NewGroupListLogic(ctx, s.svcCtx)
	return l.GroupList(in)
}

func (s *SocialServer) GroupMemberList(ctx context.Context, in *rpc.GroupMemberListReq) (*rpc.GroupMemberListResp, error) {
	l := logic.NewGroupMemberListLogic(ctx, s.svcCtx)
	return l.GroupMemberList(in)
}
