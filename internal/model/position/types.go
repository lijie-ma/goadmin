package position

import "goadmin/internal/model/schema"

// ListRequest 位置列表请求
type ListRequest struct {
	schema.PageRequest
	Keyword string `form:"keyword"` // 搜索关键词（城市或位置）
	City    string `form:"city"`    // 城市筛选
}

// CreatePositionRequest 创建位置请求
type CreatePositionRequest struct {
	City       string  `json:"city"`
	Location   string  `json:"location" binding:"required,max=128"`
	Longitude  float64 `json:"longitude" binding:"required"`
	Latitude   float64 `json:"latitude" binding:"required"`
	CustomName string  `json:"custom_name" binding:"max=128"`
}

// UpdatePositionRequest 更新位置请求
type UpdatePositionRequest struct {
	ID uint64 `json:"id" binding:"required"`
	CreatePositionRequest
}
