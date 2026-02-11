package server

import (
	"context"
	"errors"
	"goadmin/internal/model/server"
	"goadmin/pkg/db"

	"gorm.io/gorm"
)

// 确保ServerSettingRepositoryImpl实现了ServerSettingRepository接口
var _ ServerSettingRepository = (*ServerSettingRepositoryImpl)(nil)

// ServerSettingRepositoryImpl 实现ServerSettingRepository接口
type ServerSettingRepositoryImpl struct {
	*db.BaseRepository[server.ServerSetting]
}

// NewServerSettingRepository 创建服务端配置仓储实例（Wire 注入）
func NewServerSettingRepository(dbInstance *gorm.DB) ServerSettingRepository {
	return &ServerSettingRepositoryImpl{
		db.NewBaseRepository[server.ServerSetting](dbInstance),
	}
}

// Deprecated: NewServerSettingRepository_legacy 使用全局db（兼容旧代码）
func NewServerSettingRepository_legacy() ServerSettingRepository {
	return NewServerSettingRepository(db.GetDB())
}

// GetByName 根据名称获取服务端配置
func (r *ServerSettingRepositoryImpl) GetByName(ctx context.Context, name string) (*server.ServerSetting, error) {
	var setting server.ServerSetting
	err := r.DB().WithContext(ctx).Where("name = ?", name).First(&setting).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &setting, nil
}

// Exists 检查服务端配置是否存在
func (r *ServerSettingRepositoryImpl) Exists(ctx context.Context, name string) (bool, error) {
	var count int64
	err := r.DB().WithContext(ctx).Model(&server.ServerSetting{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

// BatchGet 批量获取服务端配置
func (r *ServerSettingRepositoryImpl) BatchGet(ctx context.Context, names []string) ([]*server.ServerSetting, error) {
	var settings []*server.ServerSetting
	if len(names) == 0 {
		return settings, nil
	}

	err := r.DB().WithContext(ctx).Where("name IN ?", names).Find(&settings).Error
	return settings, err
}
