// Code generated by goctl. DO NOT EDIT.
// goctl 1.8.3
// Source: group.proto

package server

import (
	"context"

	"IM/rpc/group/group"
	"IM/rpc/group/internal/logic"
	"IM/rpc/group/internal/svc"
)

type GroupServiceServer struct {
	svcCtx *svc.ServiceContext
	group.UnimplementedGroupServiceServer
}

func NewGroupServiceServer(svcCtx *svc.ServiceContext) *GroupServiceServer {
	return &GroupServiceServer{
		svcCtx: svcCtx,
	}
}

// 创建群组
func (s *GroupServiceServer) CreateGroup(ctx context.Context, in *group.CreateGroupRequest) (*group.CreateGroupResponse, error) {
	l := logic.NewCreateGroupLogic(ctx, s.svcCtx)
	return l.CreateGroup(in)
}

// 申请加入群组
func (s *GroupServiceServer) JoinGroup(ctx context.Context, in *group.JoinGroupRequest) (*group.JoinGroupResponse, error) {
	l := logic.NewJoinGroupLogic(ctx, s.svcCtx)
	return l.JoinGroup(in)
}

// 获取加入群组申请列表
func (s *GroupServiceServer) GetJoinGroupApplications(ctx context.Context, in *group.GetJoinGroupApplicationsRequest) (*group.GetJoinGroupApplicationsResponse, error) {
	l := logic.NewGetJoinGroupApplicationsLogic(ctx, s.svcCtx)
	return l.GetJoinGroupApplications(in)
}

// 处理加入群组申请
func (s *GroupServiceServer) HandleJoinGroupApplication(ctx context.Context, in *group.HandleJoinGroupApplicationRequest) (*group.HandleJoinGroupApplicationResponse, error) {
	l := logic.NewHandleJoinGroupApplicationLogic(ctx, s.svcCtx)
	return l.HandleJoinGroupApplication(in)
}

// 搜索群组
func (s *GroupServiceServer) SearchGroup(ctx context.Context, in *group.SearchGroupRequest) (*group.SearchGroupResponse, error) {
	l := logic.NewSearchGroupLogic(ctx, s.svcCtx)
	return l.SearchGroup(in)
}

// 邀请加入群组
func (s *GroupServiceServer) InviteToGroup(ctx context.Context, in *group.InviteToGroupRequest) (*group.InviteToGroupResponse, error) {
	l := logic.NewInviteToGroupLogic(ctx, s.svcCtx)
	return l.InviteToGroup(in)
}

// 退出群组
func (s *GroupServiceServer) LeaveGroup(ctx context.Context, in *group.LeaveGroupRequest) (*group.LeaveGroupResponse, error) {
	l := logic.NewLeaveGroupLogic(ctx, s.svcCtx)
	return l.LeaveGroup(in)
}

// 踢出群组
func (s *GroupServiceServer) KickFromGroup(ctx context.Context, in *group.KickFromGroupRequest) (*group.KickFromGroupResponse, error) {
	l := logic.NewKickFromGroupLogic(ctx, s.svcCtx)
	return l.KickFromGroup(in)
}

// 获取群组信息
func (s *GroupServiceServer) GetGroupInfo(ctx context.Context, in *group.GetGroupInfoRequest) (*group.GetGroupInfoResponse, error) {
	l := logic.NewGetGroupInfoLogic(ctx, s.svcCtx)
	return l.GetGroupInfo(in)
}

// 获取群组列表
func (s *GroupServiceServer) GetGroupList(ctx context.Context, in *group.GetGroupListRequest) (*group.GetGroupListResponse, error) {
	l := logic.NewGetGroupListLogic(ctx, s.svcCtx)
	return l.GetGroupList(in)
}

// 获取群组成员列表
func (s *GroupServiceServer) GetGroupMemberList(ctx context.Context, in *group.GetGroupMemberListRequest) (*group.GetGroupMemberListResponse, error) {
	l := logic.NewGetGroupMemberListLogic(ctx, s.svcCtx)
	return l.GetGroupMemberList(in)
}

// 更新群组信息
func (s *GroupServiceServer) UpdateGroupInfo(ctx context.Context, in *group.UpdateGroupInfoRequest) (*group.UpdateGroupInfoResponse, error) {
	l := logic.NewUpdateGroupInfoLogic(ctx, s.svcCtx)
	return l.UpdateGroupInfo(in)
}

// 设置群组成员角色
func (s *GroupServiceServer) SetMemberRole(ctx context.Context, in *group.SetMemberRoleRequest) (*group.SetMemberRoleResponse, error) {
	l := logic.NewSetMemberRoleLogic(ctx, s.svcCtx)
	return l.SetMemberRole(in)
}

// 禁言群组成员
func (s *GroupServiceServer) MuteMember(ctx context.Context, in *group.MuteMemberRequest) (*group.MuteMemberResponse, error) {
	l := logic.NewMuteMemberLogic(ctx, s.svcCtx)
	return l.MuteMember(in)
}

// 解散群组
func (s *GroupServiceServer) DismissGroup(ctx context.Context, in *group.DismissGroupRequest) (*group.DismissGroupResponse, error) {
	l := logic.NewDismissGroupLogic(ctx, s.svcCtx)
	return l.DismissGroup(in)
}

// 转让群组
func (s *GroupServiceServer) TransferGroup(ctx context.Context, in *group.TransferGroupRequest) (*group.TransferGroupResponse, error) {
	l := logic.NewTransferGroupLogic(ctx, s.svcCtx)
	return l.TransferGroup(in)
}

// 获取群组成员信息
func (s *GroupServiceServer) GetGroupMemberInfo(ctx context.Context, in *group.GetGroupMemberInfoRequest) (*group.GetGroupMemberInfoResponse, error) {
	l := logic.NewGetGroupMemberInfoLogic(ctx, s.svcCtx)
	return l.GetGroupMemberInfo(in)
}

// 修改群组成员信息
func (s *GroupServiceServer) UpdateGroupMemberInfo(ctx context.Context, in *group.UpdateGroupMemberInfoRequest) (*group.UpdateGroupMemberInfoResponse, error) {
	l := logic.NewUpdateGroupMemberInfoLogic(ctx, s.svcCtx)
	return l.UpdateGroupMemberInfo(in)
}

// 获取群组通知列表 (调用后，返回的通知在后端被标记为已读)
func (s *GroupServiceServer) GetGroupNotifications(ctx context.Context, in *group.GetGroupNotificationsRequest) (*group.GetGroupNotificationsResponse, error) {
	l := logic.NewGetGroupNotificationsLogic(ctx, s.svcCtx)
	return l.GetGroupNotifications(in)
}

// 获取总的未读数
func (s *GroupServiceServer) GetUnreadCount(ctx context.Context, in *group.GetUnreadCountRequest) (*group.GetUnreadCountResponse, error) {
	l := logic.NewGetUnreadCountLogic(ctx, s.svcCtx)
	return l.GetUnreadCount(in)
}
