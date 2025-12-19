package role

import "goadmin/internal/model/schema"

// LoginRequest 登录请求参数
type CreateRequest struct {
	Code        string `json:"code" form:"code"`
	Name        string `json:"name" form:"name" binding:"required,max=50"`
	Description string `json:"description" form:"description" binding:"max=200"`
	Status      int    `json:"status" form:"status" binding:"oneof=0 1"`
}

type UpdateRequest struct {
	schema.IDRequest
	CreateRequest
}
