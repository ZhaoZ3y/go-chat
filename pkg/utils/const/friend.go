package _const

// 好友关系状态
const (
	FriendStatusNormal    = 1 // 正常
	FriendStatusBlocked   = 2 // 已拉黑 (我拉黑了对方)
	FriendStatusBeBlocked = 3 // 被对方拉黑
)

// 好友申请状态
const (
	FriendRequestStatusPending  = 1 // 待处理
	FriendRequestStatusAccepted = 2 // 已同意
	FriendRequestStatusRejected = 3 // 已拒绝
)
