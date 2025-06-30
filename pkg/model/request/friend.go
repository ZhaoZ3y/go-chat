package request

// SendFriendRequestReq - 发送好友申请的请求
type SendFriendRequestReq struct {
	ToUserID int64  `json:"to_user_id" binding:"required"`
	Message  string `json:"message"`
}

// HandleFriendRequestReq - 处理好友申请的请求
type HandleFriendRequestReq struct {
	RequestID int64  `json:"request_id" binding:"required"`
	Action    int32  `json:"action" binding:"required,oneof=2 3"` // 2=同意, 3=拒绝
	Message   string `json:"message"`                             // 可选的处理消息
}

// GetFriendRequestListReq - 获取好友申请列表的请求 (支持分页和状态过滤)
type GetFriendRequestListReq struct {
	Page     int   `form:"page"`      // 页码
	PageSize int   `form:"page_size"` // 每页数量
	Status   int32 `form:"status"`    // 申请状态 (0=全部, 1=待处理, 2=已同意, 3=已拒绝)
}

// FriendDeleteActionReq - 用于删除好友的通用请求
type FriendDeleteActionReq struct {
	FriendID int64 `json:"friend_id" binding:"required"`
}

// BlockFriendReq - 用于拉黑好友的请求
type BlockFriendReq struct {
	FriendID int64 `json:"friend_id" binding:"required"`
	Status   int32 `json:"status" binding:"required,oneof=0 1"` // 0=取消拉黑, 1=拉黑
}

// UpdateFriendRemarkReq - 更新好友备注的请求
type UpdateFriendRemarkReq struct {
	FriendID int64  `json:"friend_id" binding:"required"`
	Remark   string `json:"remark"`
}

// GetListReq - 通用的分页列表请求
type GetListReq struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}
