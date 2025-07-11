// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.19.4
// source: group.proto

package group

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	GroupService_CreateGroup_FullMethodName                = "/group.GroupService/CreateGroup"
	GroupService_JoinGroup_FullMethodName                  = "/group.GroupService/JoinGroup"
	GroupService_GetJoinGroupApplications_FullMethodName   = "/group.GroupService/GetJoinGroupApplications"
	GroupService_HandleJoinGroupApplication_FullMethodName = "/group.GroupService/HandleJoinGroupApplication"
	GroupService_SearchGroup_FullMethodName                = "/group.GroupService/SearchGroup"
	GroupService_InviteToGroup_FullMethodName              = "/group.GroupService/InviteToGroup"
	GroupService_LeaveGroup_FullMethodName                 = "/group.GroupService/LeaveGroup"
	GroupService_KickFromGroup_FullMethodName              = "/group.GroupService/KickFromGroup"
	GroupService_GetGroupInfo_FullMethodName               = "/group.GroupService/GetGroupInfo"
	GroupService_GetGroupList_FullMethodName               = "/group.GroupService/GetGroupList"
	GroupService_GetGroupMemberList_FullMethodName         = "/group.GroupService/GetGroupMemberList"
	GroupService_UpdateGroupInfo_FullMethodName            = "/group.GroupService/UpdateGroupInfo"
	GroupService_SetMemberRole_FullMethodName              = "/group.GroupService/SetMemberRole"
	GroupService_MuteMember_FullMethodName                 = "/group.GroupService/MuteMember"
	GroupService_DismissGroup_FullMethodName               = "/group.GroupService/DismissGroup"
	GroupService_TransferGroup_FullMethodName              = "/group.GroupService/TransferGroup"
	GroupService_GetGroupMemberInfo_FullMethodName         = "/group.GroupService/GetGroupMemberInfo"
	GroupService_UpdateGroupMemberInfo_FullMethodName      = "/group.GroupService/UpdateGroupMemberInfo"
	GroupService_GetGroupNotifications_FullMethodName      = "/group.GroupService/GetGroupNotifications"
	GroupService_GetUnreadCount_FullMethodName             = "/group.GroupService/GetUnreadCount"
)

// GroupServiceClient is the client API for GroupService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// 群组服务
type GroupServiceClient interface {
	// 创建群组
	CreateGroup(ctx context.Context, in *CreateGroupRequest, opts ...grpc.CallOption) (*CreateGroupResponse, error)
	// 申请加入群组
	JoinGroup(ctx context.Context, in *JoinGroupRequest, opts ...grpc.CallOption) (*JoinGroupResponse, error)
	// 获取加入群组申请列表
	GetJoinGroupApplications(ctx context.Context, in *GetJoinGroupApplicationsRequest, opts ...grpc.CallOption) (*GetJoinGroupApplicationsResponse, error)
	// 处理加入群组申请
	HandleJoinGroupApplication(ctx context.Context, in *HandleJoinGroupApplicationRequest, opts ...grpc.CallOption) (*HandleJoinGroupApplicationResponse, error)
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
	// 获取群组成员信息
	GetGroupMemberInfo(ctx context.Context, in *GetGroupMemberInfoRequest, opts ...grpc.CallOption) (*GetGroupMemberInfoResponse, error)
	// 修改群组成员信息
	UpdateGroupMemberInfo(ctx context.Context, in *UpdateGroupMemberInfoRequest, opts ...grpc.CallOption) (*UpdateGroupMemberInfoResponse, error)
	// 获取群组通知列表 (调用后，返回的通知在后端被标记为已读)
	GetGroupNotifications(ctx context.Context, in *GetGroupNotificationsRequest, opts ...grpc.CallOption) (*GetGroupNotificationsResponse, error)
	// 获取总的未读数
	GetUnreadCount(ctx context.Context, in *GetUnreadCountRequest, opts ...grpc.CallOption) (*GetUnreadCountResponse, error)
}

type groupServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGroupServiceClient(cc grpc.ClientConnInterface) GroupServiceClient {
	return &groupServiceClient{cc}
}

func (c *groupServiceClient) CreateGroup(ctx context.Context, in *CreateGroupRequest, opts ...grpc.CallOption) (*CreateGroupResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateGroupResponse)
	err := c.cc.Invoke(ctx, GroupService_CreateGroup_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) JoinGroup(ctx context.Context, in *JoinGroupRequest, opts ...grpc.CallOption) (*JoinGroupResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(JoinGroupResponse)
	err := c.cc.Invoke(ctx, GroupService_JoinGroup_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) GetJoinGroupApplications(ctx context.Context, in *GetJoinGroupApplicationsRequest, opts ...grpc.CallOption) (*GetJoinGroupApplicationsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetJoinGroupApplicationsResponse)
	err := c.cc.Invoke(ctx, GroupService_GetJoinGroupApplications_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) HandleJoinGroupApplication(ctx context.Context, in *HandleJoinGroupApplicationRequest, opts ...grpc.CallOption) (*HandleJoinGroupApplicationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HandleJoinGroupApplicationResponse)
	err := c.cc.Invoke(ctx, GroupService_HandleJoinGroupApplication_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) SearchGroup(ctx context.Context, in *SearchGroupRequest, opts ...grpc.CallOption) (*SearchGroupResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SearchGroupResponse)
	err := c.cc.Invoke(ctx, GroupService_SearchGroup_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) InviteToGroup(ctx context.Context, in *InviteToGroupRequest, opts ...grpc.CallOption) (*InviteToGroupResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(InviteToGroupResponse)
	err := c.cc.Invoke(ctx, GroupService_InviteToGroup_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) LeaveGroup(ctx context.Context, in *LeaveGroupRequest, opts ...grpc.CallOption) (*LeaveGroupResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LeaveGroupResponse)
	err := c.cc.Invoke(ctx, GroupService_LeaveGroup_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) KickFromGroup(ctx context.Context, in *KickFromGroupRequest, opts ...grpc.CallOption) (*KickFromGroupResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(KickFromGroupResponse)
	err := c.cc.Invoke(ctx, GroupService_KickFromGroup_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) GetGroupInfo(ctx context.Context, in *GetGroupInfoRequest, opts ...grpc.CallOption) (*GetGroupInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetGroupInfoResponse)
	err := c.cc.Invoke(ctx, GroupService_GetGroupInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) GetGroupList(ctx context.Context, in *GetGroupListRequest, opts ...grpc.CallOption) (*GetGroupListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetGroupListResponse)
	err := c.cc.Invoke(ctx, GroupService_GetGroupList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) GetGroupMemberList(ctx context.Context, in *GetGroupMemberListRequest, opts ...grpc.CallOption) (*GetGroupMemberListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetGroupMemberListResponse)
	err := c.cc.Invoke(ctx, GroupService_GetGroupMemberList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) UpdateGroupInfo(ctx context.Context, in *UpdateGroupInfoRequest, opts ...grpc.CallOption) (*UpdateGroupInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateGroupInfoResponse)
	err := c.cc.Invoke(ctx, GroupService_UpdateGroupInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) SetMemberRole(ctx context.Context, in *SetMemberRoleRequest, opts ...grpc.CallOption) (*SetMemberRoleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SetMemberRoleResponse)
	err := c.cc.Invoke(ctx, GroupService_SetMemberRole_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) MuteMember(ctx context.Context, in *MuteMemberRequest, opts ...grpc.CallOption) (*MuteMemberResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MuteMemberResponse)
	err := c.cc.Invoke(ctx, GroupService_MuteMember_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) DismissGroup(ctx context.Context, in *DismissGroupRequest, opts ...grpc.CallOption) (*DismissGroupResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DismissGroupResponse)
	err := c.cc.Invoke(ctx, GroupService_DismissGroup_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) TransferGroup(ctx context.Context, in *TransferGroupRequest, opts ...grpc.CallOption) (*TransferGroupResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TransferGroupResponse)
	err := c.cc.Invoke(ctx, GroupService_TransferGroup_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) GetGroupMemberInfo(ctx context.Context, in *GetGroupMemberInfoRequest, opts ...grpc.CallOption) (*GetGroupMemberInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetGroupMemberInfoResponse)
	err := c.cc.Invoke(ctx, GroupService_GetGroupMemberInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) UpdateGroupMemberInfo(ctx context.Context, in *UpdateGroupMemberInfoRequest, opts ...grpc.CallOption) (*UpdateGroupMemberInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateGroupMemberInfoResponse)
	err := c.cc.Invoke(ctx, GroupService_UpdateGroupMemberInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) GetGroupNotifications(ctx context.Context, in *GetGroupNotificationsRequest, opts ...grpc.CallOption) (*GetGroupNotificationsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetGroupNotificationsResponse)
	err := c.cc.Invoke(ctx, GroupService_GetGroupNotifications_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *groupServiceClient) GetUnreadCount(ctx context.Context, in *GetUnreadCountRequest, opts ...grpc.CallOption) (*GetUnreadCountResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUnreadCountResponse)
	err := c.cc.Invoke(ctx, GroupService_GetUnreadCount_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GroupServiceServer is the server API for GroupService service.
// All implementations must embed UnimplementedGroupServiceServer
// for forward compatibility.
//
// 群组服务
type GroupServiceServer interface {
	// 创建群组
	CreateGroup(context.Context, *CreateGroupRequest) (*CreateGroupResponse, error)
	// 申请加入群组
	JoinGroup(context.Context, *JoinGroupRequest) (*JoinGroupResponse, error)
	// 获取加入群组申请列表
	GetJoinGroupApplications(context.Context, *GetJoinGroupApplicationsRequest) (*GetJoinGroupApplicationsResponse, error)
	// 处理加入群组申请
	HandleJoinGroupApplication(context.Context, *HandleJoinGroupApplicationRequest) (*HandleJoinGroupApplicationResponse, error)
	// 搜索群组
	SearchGroup(context.Context, *SearchGroupRequest) (*SearchGroupResponse, error)
	// 邀请加入群组
	InviteToGroup(context.Context, *InviteToGroupRequest) (*InviteToGroupResponse, error)
	// 退出群组
	LeaveGroup(context.Context, *LeaveGroupRequest) (*LeaveGroupResponse, error)
	// 踢出群组
	KickFromGroup(context.Context, *KickFromGroupRequest) (*KickFromGroupResponse, error)
	// 获取群组信息
	GetGroupInfo(context.Context, *GetGroupInfoRequest) (*GetGroupInfoResponse, error)
	// 获取群组列表
	GetGroupList(context.Context, *GetGroupListRequest) (*GetGroupListResponse, error)
	// 获取群组成员列表
	GetGroupMemberList(context.Context, *GetGroupMemberListRequest) (*GetGroupMemberListResponse, error)
	// 更新群组信息
	UpdateGroupInfo(context.Context, *UpdateGroupInfoRequest) (*UpdateGroupInfoResponse, error)
	// 设置群组成员角色
	SetMemberRole(context.Context, *SetMemberRoleRequest) (*SetMemberRoleResponse, error)
	// 禁言群组成员
	MuteMember(context.Context, *MuteMemberRequest) (*MuteMemberResponse, error)
	// 解散群组
	DismissGroup(context.Context, *DismissGroupRequest) (*DismissGroupResponse, error)
	// 转让群组
	TransferGroup(context.Context, *TransferGroupRequest) (*TransferGroupResponse, error)
	// 获取群组成员信息
	GetGroupMemberInfo(context.Context, *GetGroupMemberInfoRequest) (*GetGroupMemberInfoResponse, error)
	// 修改群组成员信息
	UpdateGroupMemberInfo(context.Context, *UpdateGroupMemberInfoRequest) (*UpdateGroupMemberInfoResponse, error)
	// 获取群组通知列表 (调用后，返回的通知在后端被标记为已读)
	GetGroupNotifications(context.Context, *GetGroupNotificationsRequest) (*GetGroupNotificationsResponse, error)
	// 获取总的未读数
	GetUnreadCount(context.Context, *GetUnreadCountRequest) (*GetUnreadCountResponse, error)
	mustEmbedUnimplementedGroupServiceServer()
}

// UnimplementedGroupServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedGroupServiceServer struct{}

func (UnimplementedGroupServiceServer) CreateGroup(context.Context, *CreateGroupRequest) (*CreateGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateGroup not implemented")
}
func (UnimplementedGroupServiceServer) JoinGroup(context.Context, *JoinGroupRequest) (*JoinGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JoinGroup not implemented")
}
func (UnimplementedGroupServiceServer) GetJoinGroupApplications(context.Context, *GetJoinGroupApplicationsRequest) (*GetJoinGroupApplicationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetJoinGroupApplications not implemented")
}
func (UnimplementedGroupServiceServer) HandleJoinGroupApplication(context.Context, *HandleJoinGroupApplicationRequest) (*HandleJoinGroupApplicationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleJoinGroupApplication not implemented")
}
func (UnimplementedGroupServiceServer) SearchGroup(context.Context, *SearchGroupRequest) (*SearchGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchGroup not implemented")
}
func (UnimplementedGroupServiceServer) InviteToGroup(context.Context, *InviteToGroupRequest) (*InviteToGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InviteToGroup not implemented")
}
func (UnimplementedGroupServiceServer) LeaveGroup(context.Context, *LeaveGroupRequest) (*LeaveGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LeaveGroup not implemented")
}
func (UnimplementedGroupServiceServer) KickFromGroup(context.Context, *KickFromGroupRequest) (*KickFromGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method KickFromGroup not implemented")
}
func (UnimplementedGroupServiceServer) GetGroupInfo(context.Context, *GetGroupInfoRequest) (*GetGroupInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGroupInfo not implemented")
}
func (UnimplementedGroupServiceServer) GetGroupList(context.Context, *GetGroupListRequest) (*GetGroupListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGroupList not implemented")
}
func (UnimplementedGroupServiceServer) GetGroupMemberList(context.Context, *GetGroupMemberListRequest) (*GetGroupMemberListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGroupMemberList not implemented")
}
func (UnimplementedGroupServiceServer) UpdateGroupInfo(context.Context, *UpdateGroupInfoRequest) (*UpdateGroupInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateGroupInfo not implemented")
}
func (UnimplementedGroupServiceServer) SetMemberRole(context.Context, *SetMemberRoleRequest) (*SetMemberRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetMemberRole not implemented")
}
func (UnimplementedGroupServiceServer) MuteMember(context.Context, *MuteMemberRequest) (*MuteMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MuteMember not implemented")
}
func (UnimplementedGroupServiceServer) DismissGroup(context.Context, *DismissGroupRequest) (*DismissGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DismissGroup not implemented")
}
func (UnimplementedGroupServiceServer) TransferGroup(context.Context, *TransferGroupRequest) (*TransferGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TransferGroup not implemented")
}
func (UnimplementedGroupServiceServer) GetGroupMemberInfo(context.Context, *GetGroupMemberInfoRequest) (*GetGroupMemberInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGroupMemberInfo not implemented")
}
func (UnimplementedGroupServiceServer) UpdateGroupMemberInfo(context.Context, *UpdateGroupMemberInfoRequest) (*UpdateGroupMemberInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateGroupMemberInfo not implemented")
}
func (UnimplementedGroupServiceServer) GetGroupNotifications(context.Context, *GetGroupNotificationsRequest) (*GetGroupNotificationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGroupNotifications not implemented")
}
func (UnimplementedGroupServiceServer) GetUnreadCount(context.Context, *GetUnreadCountRequest) (*GetUnreadCountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUnreadCount not implemented")
}
func (UnimplementedGroupServiceServer) mustEmbedUnimplementedGroupServiceServer() {}
func (UnimplementedGroupServiceServer) testEmbeddedByValue()                      {}

// UnsafeGroupServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GroupServiceServer will
// result in compilation errors.
type UnsafeGroupServiceServer interface {
	mustEmbedUnimplementedGroupServiceServer()
}

func RegisterGroupServiceServer(s grpc.ServiceRegistrar, srv GroupServiceServer) {
	// If the following call pancis, it indicates UnimplementedGroupServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&GroupService_ServiceDesc, srv)
}

func _GroupService_CreateGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).CreateGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GroupService_CreateGroup_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).CreateGroup(ctx, req.(*CreateGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_JoinGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).JoinGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GroupService_JoinGroup_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).JoinGroup(ctx, req.(*JoinGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_GetJoinGroupApplications_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetJoinGroupApplicationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).GetJoinGroupApplications(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GroupService_GetJoinGroupApplications_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).GetJoinGroupApplications(ctx, req.(*GetJoinGroupApplicationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_HandleJoinGroupApplication_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HandleJoinGroupApplicationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).HandleJoinGroupApplication(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GroupService_HandleJoinGroupApplication_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).HandleJoinGroupApplication(ctx, req.(*HandleJoinGroupApplicationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_SearchGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).SearchGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GroupService_SearchGroup_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).SearchGroup(ctx, req.(*SearchGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_InviteToGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InviteToGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).InviteToGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GroupService_InviteToGroup_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).InviteToGroup(ctx, req.(*InviteToGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_LeaveGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaveGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).LeaveGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GroupService_LeaveGroup_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).LeaveGroup(ctx, req.(*LeaveGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_KickFromGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KickFromGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).KickFromGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GroupService_KickFromGroup_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).KickFromGroup(ctx, req.(*KickFromGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_GetGroupInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGroupInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).GetGroupInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GroupService_GetGroupInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).GetGroupInfo(ctx, req.(*GetGroupInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_GetGroupList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGroupListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).GetGroupList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GroupService_GetGroupList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).GetGroupList(ctx, req.(*GetGroupListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_GetGroupMemberList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGroupMemberListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).GetGroupMemberList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GroupService_GetGroupMemberList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).GetGroupMemberList(ctx, req.(*GetGroupMemberListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_UpdateGroupInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateGroupInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).UpdateGroupInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GroupService_UpdateGroupInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).UpdateGroupInfo(ctx, req.(*UpdateGroupInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_SetMemberRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetMemberRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).SetMemberRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GroupService_SetMemberRole_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).SetMemberRole(ctx, req.(*SetMemberRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_MuteMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MuteMemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).MuteMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GroupService_MuteMember_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).MuteMember(ctx, req.(*MuteMemberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_DismissGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DismissGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).DismissGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GroupService_DismissGroup_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).DismissGroup(ctx, req.(*DismissGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_TransferGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TransferGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).TransferGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GroupService_TransferGroup_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).TransferGroup(ctx, req.(*TransferGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_GetGroupMemberInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGroupMemberInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).GetGroupMemberInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GroupService_GetGroupMemberInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).GetGroupMemberInfo(ctx, req.(*GetGroupMemberInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_UpdateGroupMemberInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateGroupMemberInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).UpdateGroupMemberInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GroupService_UpdateGroupMemberInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).UpdateGroupMemberInfo(ctx, req.(*UpdateGroupMemberInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_GetGroupNotifications_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetGroupNotificationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).GetGroupNotifications(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GroupService_GetGroupNotifications_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).GetGroupNotifications(ctx, req.(*GetGroupNotificationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GroupService_GetUnreadCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUnreadCountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GroupServiceServer).GetUnreadCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GroupService_GetUnreadCount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GroupServiceServer).GetUnreadCount(ctx, req.(*GetUnreadCountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GroupService_ServiceDesc is the grpc.ServiceDesc for GroupService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GroupService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "group.GroupService",
	HandlerType: (*GroupServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateGroup",
			Handler:    _GroupService_CreateGroup_Handler,
		},
		{
			MethodName: "JoinGroup",
			Handler:    _GroupService_JoinGroup_Handler,
		},
		{
			MethodName: "GetJoinGroupApplications",
			Handler:    _GroupService_GetJoinGroupApplications_Handler,
		},
		{
			MethodName: "HandleJoinGroupApplication",
			Handler:    _GroupService_HandleJoinGroupApplication_Handler,
		},
		{
			MethodName: "SearchGroup",
			Handler:    _GroupService_SearchGroup_Handler,
		},
		{
			MethodName: "InviteToGroup",
			Handler:    _GroupService_InviteToGroup_Handler,
		},
		{
			MethodName: "LeaveGroup",
			Handler:    _GroupService_LeaveGroup_Handler,
		},
		{
			MethodName: "KickFromGroup",
			Handler:    _GroupService_KickFromGroup_Handler,
		},
		{
			MethodName: "GetGroupInfo",
			Handler:    _GroupService_GetGroupInfo_Handler,
		},
		{
			MethodName: "GetGroupList",
			Handler:    _GroupService_GetGroupList_Handler,
		},
		{
			MethodName: "GetGroupMemberList",
			Handler:    _GroupService_GetGroupMemberList_Handler,
		},
		{
			MethodName: "UpdateGroupInfo",
			Handler:    _GroupService_UpdateGroupInfo_Handler,
		},
		{
			MethodName: "SetMemberRole",
			Handler:    _GroupService_SetMemberRole_Handler,
		},
		{
			MethodName: "MuteMember",
			Handler:    _GroupService_MuteMember_Handler,
		},
		{
			MethodName: "DismissGroup",
			Handler:    _GroupService_DismissGroup_Handler,
		},
		{
			MethodName: "TransferGroup",
			Handler:    _GroupService_TransferGroup_Handler,
		},
		{
			MethodName: "GetGroupMemberInfo",
			Handler:    _GroupService_GetGroupMemberInfo_Handler,
		},
		{
			MethodName: "UpdateGroupMemberInfo",
			Handler:    _GroupService_UpdateGroupMemberInfo_Handler,
		},
		{
			MethodName: "GetGroupNotifications",
			Handler:    _GroupService_GetGroupNotifications_Handler,
		},
		{
			MethodName: "GetUnreadCount",
			Handler:    _GroupService_GetUnreadCount_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "group.proto",
}
