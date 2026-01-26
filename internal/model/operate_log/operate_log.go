package operate_log

import (
	"goadmin/internal/model/schema"
)

// OperateLog 操作日志表
type OperateLog struct {
	schema.BaseModel
	Content  string `gorm:"size:512;comment:详情内容" json:"content"`
	Username string `gorm:"size:64;not null;default:'';comment:操作用户" json:"username"`
	IP       string `gorm:"size:40;not null;default:'';comment:操作人ip" json:"ip"`
}

// TableName 指定表名
func (OperateLog) TableName() string {
	return "operate_log"
}

// ListRequest 操作日志列表请求
type ListRequest struct {
	schema.PageRequest
	Username  string `form:"username" json:"username"`     // 用户名
	Content   string `form:"content" json:"content"`       // 内容
	IP        string `form:"ip" json:"ip"`                 // IP地址
	StartTime string `form:"start_time" json:"start_time"` // 开始时间
	EndTime   string `form:"end_time" json:"end_time"`     // 结束时间
}
