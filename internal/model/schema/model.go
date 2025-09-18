package schema

import "time"

// BaseModel 基础模型，包含共有字段
type BaseModel struct {
	ID    uint64    `gorm:"primary_key;auto_increment;default:0" json:"id"`
	CTime time.Time `gorm:"column:ctime;not null;default:CURRENT_TIMESTAMP" json:"ctime"`
	MTime time.Time `gorm:"column:mtime;not null;default:CURRENT_TIMESTAMP" json:"mtime"`
}

// BeforeCreate 在创建记录前自动设置时间
func (m *BaseModel) BeforeCreate() error {
	now := time.Now()
	m.CTime = now
	m.MTime = now
	return nil
}

// BeforeUpdate 在更新记录前自动设置修改时间
func (m *BaseModel) BeforeUpdate() error {
	m.MTime = time.Now()
	return nil
}
