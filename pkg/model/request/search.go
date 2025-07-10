package request

// SearchRequest 搜索用户请求
type SearchRequest struct {
	Keyword string `form:"keyword" binding:"required"`
}
