// Code generated by goctl. DO NOT EDIT.
// goctl 1.8.3
// Source: group.proto

package groupservice

import (
	"context"

	"IM/rpc/group/group"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	CreateGroupRequest         = group.CreateGroupRequest
	CreateGroupResponse        = group.CreateGroupResponse
	DismissGroupRequest        = group.DismissGroupRequest
	DismissGroupResponse       = group.DismissGroupResponse
	GetGroupInfoRequest        = group.GetGroupInfoRequest
	GetGroupInfoResponse       = group.GetGroupInfoResponse
	GetGroupListRequest        = group.GetGroupListRequest
	GetGroupListResponse       = group.GetGroupListResponse
	GetGroupMemberListRequest  = group.GetGroupMemberListRequest
	GetGroupMemberListResponse = group.GetGroupMemberListResponse
	Group                      = group.Group
	GroupMember                = group.GroupMember
	InviteToGroupRequest       = group.InviteToGroupRequest
	InviteToGroupResponse      = group.InviteToGroupResponse
	JoinGroupRequest           = group.JoinGroupRequest
	JoinGroupResponse          = group.JoinGroupResponse
	KickFromGroupRequest       = group.KickFromGroupRequest
	KickFromGroupResponse      = group.KickFromGroupResponse
	LeaveGroupRequest          = group.LeaveGroupRequest
	LeaveGroupResponse         = group.LeaveGroupResponse
	MuteMemberRequest          = group.MuteMemberRequest
	MuteMemberResponse         = group.MuteMemberResponse
	SearchGroupRequest         = group.SearchGroupRequest
	SearchGroupResponse        = group.SearchGroupResponse
	SetMemberRoleRequest       = group.SetMemberRoleRequest
	SetMemberRoleResponse      = group.SetMemberRoleResponse
	TransferGroupRequest       = group.TransferGroupRequest
	TransferGroupResponse      = group.TransferGroupResponse
	UpdateGroupInfoRequest     = group.UpdateGroupInfoRequest
	UpdateGroupInfoResponse    = group.UpdateGroupInfoResponse

	GroupService interface {
		// 创建群组
		CreateGroup(ctx context.Context, in *CreateGroupRequest, opts ...grpc.CallOption) (*CreateGroupResponse, error)
		// 加入群组
		JoinGroup(ctx context.Context, in *JoinGroupRequest, opts ...grpc.CallOption) (*JoinGroupResponse, error)
		// 搜索群组
		SearchGroup(ctx context.Context, in *SearchGroupRequest, opts ...grpc.CallOption) (*SearchGroupResponse, error)
		// 邀请加入群组
		InviteToGroup(ctx context.Context, in *InviteToGroupRequest, opts ...grpc.CallOption) (*InviteToGroupResponse, error)
		// 退出群组
		LeaveGroup(ctx context.Context, in *LeaveGroupRequest, opts ...grpc.CallOption) (*LeaveGroupResponse, error)
		// 踢出群组
		KickFromGroup(ctx context.Context, in *KickFromGroupRequest, opts ...grpc.CallOption) (*KickFromGroupResponse, error)
		// 获取群组信息
		GetGroupInfo(ctx context.Context, in *GetGroupInfoRequest, opts ...grpc.CallOption) (*GetGroupInfoResponse, error)
		// 获取群组列表
		GetGroupList(ctx context.Context, in *GetGroupListRequest, opts ...grpc.CallOption) (*GetGroupListResponse, error)
		// 获取群组成员列表
		GetGroupMemberList(ctx context.Context, in *GetGroupMemberListRequest, opts ...grpc.CallOption) (*GetGroupMemberListResponse, error)
		// 更新群组信息
		UpdateGroupInfo(ctx context.Context, in *UpdateGroupInfoRequest, opts ...grpc.CallOption) (*UpdateGroupInfoResponse, error)
		// 设置群组成员角色
		SetMemberRole(ctx context.Context, in *SetMemberRoleRequest, opts ...grpc.CallOption) (*SetMemberRoleResponse, error)
		// 禁言群组成员
		MuteMember(ctx context.Context, in *MuteMemberRequest, opts ...grpc.CallOption) (*MuteMemberResponse, error)
		// 解散群组
		DismissGroup(ctx context.Context, in *DismissGroupRequest, opts ...grpc.CallOption) (*DismissGroupResponse, error)
		// 转让群组
		TransferGroup(ctx context.Context, in *TransferGroupRequest, opts ...grpc.CallOption) (*TransferGroupResponse, error)
	}

	defaultGroupService struct {
		cli zrpc.Client
	}
)

func NewGroupService(cli zrpc.Client) GroupService {
	return &defaultGroupService{
		cli: cli,
	}
}

// 创建群组
func (m *defaultGroupService) CreateGroup(ctx context.Context, in *CreateGroupRequest, opts ...grpc.CallOption) (*CreateGroupResponse, error) {
	client := group.NewGroupServiceClient(m.cli.Conn())
	return client.CreateGroup(ctx, in, opts...)
}

// 加入群组
func (m *defaultGroupService) JoinGroup(ctx context.Context, in *JoinGroupRequest, opts ...grpc.CallOption) (*JoinGroupResponse, error) {
	client := group.NewGroupServiceClient(m.cli.Conn())
	return client.JoinGroup(ctx, in, opts...)
}

// 搜索群组
func (m *defaultGroupService) SearchGroup(ctx context.Context, in *SearchGroupRequest, opts ...grpc.CallOption) (*SearchGroupResponse, error) {
	client := group.NewGroupServiceClient(m.cli.Conn())
	return client.SearchGroup(ctx, in, opts...)
}

// 邀请加入群组
func (m *defaultGroupService) InviteToGroup(ctx context.Context, in *InviteToGroupRequest, opts ...grpc.CallOption) (*InviteToGroupResponse, error) {
	client := group.NewGroupServiceClient(m.cli.Conn())
	return client.InviteToGroup(ctx, in, opts...)
}

// 退出群组
func (m *defaultGroupService) LeaveGroup(ctx context.Context, in *LeaveGroupRequest, opts ...grpc.CallOption) (*LeaveGroupResponse, error) {
	client := group.NewGroupServiceClient(m.cli.Conn())
	return client.LeaveGroup(ctx, in, opts...)
}

// 踢出群组
func (m *defaultGroupService) KickFromGroup(ctx context.Context, in *KickFromGroupRequest, opts ...grpc.CallOption) (*KickFromGroupResponse, error) {
	client := group.NewGroupServiceClient(m.cli.Conn())
	return client.KickFromGroup(ctx, in, opts...)
}

// 获取群组信息
func (m *defaultGroupService) GetGroupInfo(ctx context.Context, in *GetGroupInfoRequest, opts ...grpc.CallOption) (*GetGroupInfoResponse, error) {
	client := group.NewGroupServiceClient(m.cli.Conn())
	return client.GetGroupInfo(ctx, in, opts...)
}

// 获取群组列表
func (m *defaultGroupService) GetGroupList(ctx context.Context, in *GetGroupListRequest, opts ...grpc.CallOption) (*GetGroupListResponse, error) {
	client := group.NewGroupServiceClient(m.cli.Conn())
	return client.GetGroupList(ctx, in, opts...)
}

// 获取群组成员列表
func (m *defaultGroupService) GetGroupMemberList(ctx context.Context, in *GetGroupMemberListRequest, opts ...grpc.CallOption) (*GetGroupMemberListResponse, error) {
	client := group.NewGroupServiceClient(m.cli.Conn())
	return client.GetGroupMemberList(ctx, in, opts...)
}

// 更新群组信息
func (m *defaultGroupService) UpdateGroupInfo(ctx context.Context, in *UpdateGroupInfoRequest, opts ...grpc.CallOption) (*UpdateGroupInfoResponse, error) {
	client := group.NewGroupServiceClient(m.cli.Conn())
	return client.UpdateGroupInfo(ctx, in, opts...)
}

// 设置群组成员角色
func (m *defaultGroupService) SetMemberRole(ctx context.Context, in *SetMemberRoleRequest, opts ...grpc.CallOption) (*SetMemberRoleResponse, error) {
	client := group.NewGroupServiceClient(m.cli.Conn())
	return client.SetMemberRole(ctx, in, opts...)
}

// 禁言群组成员
func (m *defaultGroupService) MuteMember(ctx context.Context, in *MuteMemberRequest, opts ...grpc.CallOption) (*MuteMemberResponse, error) {
	client := group.NewGroupServiceClient(m.cli.Conn())
	return client.MuteMember(ctx, in, opts...)
}

// 解散群组
func (m *defaultGroupService) DismissGroup(ctx context.Context, in *DismissGroupRequest, opts ...grpc.CallOption) (*DismissGroupResponse, error) {
	client := group.NewGroupServiceClient(m.cli.Conn())
	return client.DismissGroup(ctx, in, opts...)
}

// 转让群组
func (m *defaultGroupService) TransferGroup(ctx context.Context, in *TransferGroupRequest, opts ...grpc.CallOption) (*TransferGroupResponse, error) {
	client := group.NewGroupServiceClient(m.cli.Conn())
	return client.TransferGroup(ctx, in, opts...)
}
