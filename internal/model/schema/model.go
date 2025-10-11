package schema

import (
	"goadmin/pkg/util"

	"gorm.io/gorm"
)

// BaseModel 基础模型，包含共有字段
type BaseModel struct {
	ID    uint64        `gorm:"primary_key;auto_increment;default:0" json:"id"`
	CTime util.DateTime `gorm:"column:ctime" json:"ctime"`
	MTime util.DateTime `gorm:"column:mtime" json:"mtime"`
}

// BeforeCreate 在创建记录前自动设置时间
func (m *BaseModel) BeforeCreate(tx *gorm.DB) error {
	now := util.Now()
	if m.CTime.IsZero() {
		m.CTime = now
	}
	if m.MTime.IsZero() {
		m.MTime = now
	}
	return nil
}

// BeforeUpdate 在更新记录前自动设置修改时间
func (m *BaseModel) BeforeUpdate(tx *gorm.DB) error {
	m.MTime = util.Now()
	return nil
}
