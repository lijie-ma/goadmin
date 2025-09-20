package server

import (
	"context"
	"goadmin/internal/model/server"
	"goadmin/pkg/db"
)

// ServerSettingRepository 定义服务端配置仓储接口
type ServerSettingRepository interface {
	db.Repository[server.ServerSetting]

	// GetByName 根据名称获取服务端配置
	GetByName(ctx context.Context, name string) (*server.ServerSetting, error)

	// Update 更新服务端配置
	Update(ctx context.Context, setting *server.ServerSetting) error

	// Exists 检查服务端配置是否存在
	Exists(ctx context.Context, name string) (bool, error)

	// BatchGet 批量获取服务端配置
	BatchGet(ctx context.Context, names []string) ([]*server.ServerSetting, error)
}
