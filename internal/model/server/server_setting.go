package server

import "goadmin/internal/model/schema"

// ServerSetting 服务端配置表
type ServerSetting struct {
	schema.BaseModel
	Name  string `gorm:"size:64;not null;unique;default:''" json:"name"`
	Value string `gorm:"type:text" json:"value"`
}

// TableName 指定表名
func (ServerSetting) TableName() string {
	return "server_setting"
}
