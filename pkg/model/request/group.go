package request

import "IM/rpc/group/group"

// CreateGroupRequest 创建群组请求结构
type CreateGroupRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Avatar      string  `json:"avatar"`
	MemberIds   []int64 `json:"member_ids"`
}

// UpdateGroupRequest 更新群组信息请求结构
type UpdateGroupRequest struct {
	GroupId     int64  `json:"group_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
}

// UpdateGroupMemberInfoRequest 更新群成员信息 (如群昵称)
type UpdateGroupMemberInfoRequest struct {
	GroupId  int64  `json:"group_id"`
	Nickname string `json:"nickname"`
}

// SetMemberRoleRequest 设置成员角色请求结构
type SetMemberRoleRequest struct {
	GroupId int64            `json:"group_id"`
	UserId  int64            `json:"user_id"`
	Role    group.MemberRole `json:"role"`
}

// MuteMemberRequest 禁言/解禁成员请求
type MuteMemberRequest struct {
	GroupId  int64 `json:"group_id"`
	UserId   int64 `json:"user_id"`
	Duration int64 `json:"duration"`  // 禁言时长(秒), 解禁时可不传
	IsUnmute bool  `json:"is_unmute"` // 是否解除禁言
}

// InviteToGroupRequest 邀请进群请求
type InviteToGroupRequest struct {
	GroupId int64   `json:"group_id"`
	UserIds []int64 `json:"user_ids"`
}

// JoinGroupRequest 申请加入群组请求
type JoinGroupRequest struct {
	GroupId int64  `json:"group_id"`
	Reason  string `json:"reason"`
}

// HandleJoinGroupApplicationRequest 处理入群申请请求
type HandleJoinGroupApplicationRequest struct {
	ApplicationId int64 `json:"application_id"`
	Approve       bool  `json:"approve"`
}

// KickFromGroupRequest 踢出群成员请求
type KickFromGroupRequest struct {
	GroupId int64 `json:"group_id"`
	UserId  int64 `json:"user_id"`
}

// LeaveGroupRequest 退出群组请求
type LeaveGroupRequest struct {
	GroupId int64 `json:"group_id"`
}

// DismissGroupRequest 解散群组请求
type DismissGroupRequest struct {
	GroupId int64 `json:"group_id"`
}

// TransferGroupRequest 转让群组请求
type TransferGroupRequest struct {
	GroupId    int64 `json:"group_id"`
	NewOwnerId int64 `json:"new_owner_id"`
}
