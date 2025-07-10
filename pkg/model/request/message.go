package request

// SendMessageRequest 定义了发送消息接口的请求体结构
type SendMessageRequest struct {
	ToUserId int64  `json:"to_user_id,omitempty"`       // 私聊时对方的用户ID
	GroupId  int64  `json:"group_id,omitempty"`         // 群聊时的群组ID
	Type     int64  `json:"type"`                       // 消息类型，例如 0:文本, 1:图片, 2:语音等
	Content  string `json:"content" binding:"required"` // 消息内容
	Extra    string `json:"extra,omitempty"`            // 额外信息，可用于存储JSON格式的附加数据
	ChatType int64  `json:"chat_type"`                  // 聊天类型, 0:私聊, 1:群聊
}

// GetMessageHistoryRequest 定义了获取消息历史接口的查询参数结构
type GetMessageHistoryRequest struct {
	TargetId      int64 `json:"target_id"`       // 目标ID（私聊是对方用户ID，群聊是群ID）
	ChatType      int32 `json:"chat_type"`       // 聊天类型, 0:私聊, 1:群聊
	LastMessageId int64 `json:"last_message_id"` // 上一次加载的最后一条消息ID，用于分页
	Limit         int64 `json:"limit"`           // 每页数量，默认为20
	Date          int64 `json:"date"`            // 按天查询的Unix时间戳，如果提供则忽略 LastMessageId
}

// MarkMessageReadRequest 定义了标记消息已读接口的请求体结构
type MarkMessageReadRequest struct {
	TargetId          int64 `json:"target_id"`            // 目标ID（私聊是对方用户ID，群聊是群ID）
	ChatType          int64 `json:"chat_type"`            // 聊天类型, 0:私聊, 1:群聊
	LastReadMessageId int64 `json:"last_read_message_id"` // 用户已读的最后一条消息的ID
}

// RecallMessageRequest 定义了撤回消息接口的请求体结构
type RecallMessageRequest struct {
	MessageId int64 `json:"message_id"` // 要撤回的消息ID
}

// DeleteMessageRequest 定义了删除消息接口的请求体结构
type DeleteMessageRequest struct {
	MessageId int64 `json:"message_id"` // 要从用户视图中删除的消息ID
}

// DeleteConversationRequest 定义了删除会话接口的请求体结构
type DeleteConversationRequest struct {
	TargetId int64 `json:"target_id"` // 目标ID（私聊是对方用户ID，群聊是群ID）
	ChatType int64 `json:"chat_type"` // 聊天类型, 0:私聊, 1:群聊
}

// PinConversationRequest 定义了置顶会话接口的请求体结构
type PinConversationRequest struct {
	TargetId int64 `json:"target_id"` // 目标ID（私聊是对方用户ID，群聊是群ID）
	ChatType int64 `json:"chat_type"` // 聊天类型, 0:私聊, 1:群聊
	IsPinned bool  `json:"is_pinned"` // 是否置顶
}
