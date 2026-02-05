package position

import (
	"goadmin/internal/model/schema"
)

type Position struct {
	schema.BaseModel
	City       string  `gorm:"column:city;type:varchar(64);not null;default:'';index:idx_city;comment:城市名称"`
	Location   string  `gorm:"column:location;type:varchar(128);not null;default:'';index:idx_location;comment:详细位置（如街道/建筑）"`
	Longitude  float64 `gorm:"column:longitude;type:numeric(10,6);not null;comment:经度"`
	Latitude   float64 `gorm:"column:latitude;type:numeric(10,6);not null;comment:纬度"`
	CustomName string  `gorm:"column:custom_name;type:varchar(128);comment:自定义名称"`
	CreatorID  int     `gorm:"column:creator_id;type:int;default:0;comment:创建人ID"`
	Creator    string  `gorm:"column:creator;type:varchar(64);not null;default:'';comment:创建人"`
}

// TableName 指定表名
func (Position) TableName() string {
	return "positions"
}
