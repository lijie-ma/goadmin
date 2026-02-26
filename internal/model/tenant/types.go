package tenant

import "goadmin/internal/model/schema"

// TenantStatus 租户状态
type TenantStatus int

const (
	TenantStatusEnabled  TenantStatus = 1 // 启用
	TenantStatusDisabled TenantStatus = 2 // 停用
)

// ListRequest 租户列表请求参数
type ListRequest struct {
	schema.PageRequest
	Keyword string `form:"keyword"` // 搜索关键词（名称或编码）
	Status  *int   `form:"status"`  // 状态筛选
}

// CreateRequest 创建租户请求参数
type CreateRequest struct {
	Name         string `json:"name" binding:"required,max=128"`          // 租户名称
	Code         string `json:"code" binding:"required,max=64"`           // 租户编码
	ContactEmail string `json:"contact_email" binding:"omitempty,email"`  // 联系邮箱
	ContactPhone string `json:"contact_phone" binding:"omitempty,max=32"` // 联系电话
	Status       int    `json:"status" binding:"omitempty,oneof=0 1"`     // 状态
	Config       string `json:"config"`                                   // 扩展配置JSON
}

// UpdateRequest 更新租户请求参数
type UpdateRequest struct {
	schema.IDRequest
	CreateRequest
}
