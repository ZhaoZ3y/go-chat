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
