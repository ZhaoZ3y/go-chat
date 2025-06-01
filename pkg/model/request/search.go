package request

// SearchRequest 搜索用户请求
type SearchRequest struct {
	Keyword  string `form:"keyword" binding:"required"`
	Page     int32  `form:"page"`
	PageSize int32  `form:"page_size"`
}
