package tenant

import (
	"goadmin/internal/model/schema"
)

// Tenant 租户表
type Tenant struct {
	schema.BaseModel
	Name         string       `gorm:"column:name;size:128;not null;default:''" json:"name"`
	Code         string       `gorm:"column:code;size:64;unique;not null;default:'';comment:租户唯一编码" json:"code"`
	ContactEmail string       `gorm:"column:contact_email;size:128;default:''" json:"contact_email"`
	ContactPhone string       `gorm:"column:contact_phone;size:32;default:''" json:"contact_phone"`
	Status       TenantStatus `gorm:"column:status;type:tinyint;default:1;comment:1启用 2停用" json:"status"`
	Config       string       `gorm:"column:config;type:json;comment:扩展配置，如logo、域名、自定义参数" json:"config"`
}

// TableName 指定表名
func (Tenant) TableName() string {
	return "tenants"
}

// IsActive 判断租户是否启用
func (t *Tenant) IsActive() bool {
	return t.Status == TenantStatusEnabled
}
