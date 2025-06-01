package request

// CreateGroupRequest 创建群组请求结构
type CreateGroupRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=50"`
	Description string  `json:"description" binding:"max=200"`
	Avatar      string  `json:"avatar"`
	MemberIds   []int64 `json:"member_ids"`
}

// UpdateGroupRequest 更新群组请求结构
type UpdateGroupRequest struct {
	GroupId     int64  `json:"group_id" binding:"required"`
	Name        string `json:"name" binding:"omitempty,min=1,max=50"`
	Description string `json:"description" binding:"omitempty,max=200"`
	Avatar      string `json:"avatar"`
}

// SetMemberRoleRequest 设置成员角色请求结构
type SetMemberRoleRequest struct {
	GroupId int64 `json:"group_id" binding:"required"`
	UserId  int64 `json:"user_id" binding:"required"`
	Role    int32 `json:"role" binding:"required,min=1,max=3"`
}

// MuteMemberRequest 禁言成员请求
type MuteMemberRequest struct {
	GroupId  int64 `json:"group_id" binding:"required"`
	UserId   int64 `json:"user_id" binding:"required"`
	Duration int64 `json:"duration"` // 禁言时长(秒)，0表示解除禁言
}

// InviteToGroupRequest 邀请进群请求
type InviteToGroupRequest struct {
	GroupId int64   `json:"group_id" binding:"required"`
	UserIds []int64 `json:"user_ids" binding:"required"`
}

// JoinGroupRequest 加入群组请求
type JoinGroupRequest struct {
	GroupId int64  `json:"group_id" binding:"required"`
	Reason  string `json:"reason"`
}

// KickFromGroupRequest 踢出群成员请求
type KickFromGroupRequest struct {
	GroupId int64 `json:"group_id" binding:"required"`
	UserId  int64 `json:"user_id" binding:"required"`
}

// LeaveGroupRequest 退出群组请求
type LeaveGroupRequest struct {
	GroupId int64 `json:"group_id" binding:"required"`
}

// DismissGroupRequest 解散群组请求
type DismissGroupRequest struct {
	GroupId int64 `json:"group_id" binding:"required"`
}

// TransferGroupRequest 转让群组请求
type TransferGroupRequest struct {
	GroupId    int64 `json:"group_id" binding:"required"`
	NewOwnerId int64 `json:"new_owner_id" binding:"required"`
}
