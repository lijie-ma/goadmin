package schema

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageRequest 分页请求
type PageRequest struct {
	Page     int    `form:"page,default=1" binding:"min=1"`       // 页码
	PageSize int    `form:"page_size,default=10" binding:"min=1"` // 每页数量
	OrderBy  string `form:"order_by,default=id desc"`             // 排序字段
}

// PageResponse 分页响应
type PageResponse struct {
	Total int64       `json:"total"` // 总数
	List  interface{} `json:"list"`  // 列表
}

// IDRequest ID请求
type IDRequest struct {
	ID uint64 `form:"id" json:"id" binding:"required"` // ID
}
