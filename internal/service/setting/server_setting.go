package setting

import (
	"context"
	"encoding/json"
	"goadmin/internal/model/server"
	serverRepo "goadmin/internal/repository/server"
	"goadmin/pkg/db"
)

// ServerSettingService 服务端设置服务接口
type ServerSettingService interface {
	// GetByName 根据名称获取服务端配置
	GetByName(ctx context.Context, name string) (*server.ServerSetting, error)

	// SetByName 设置服务端配置
	SetByName(ctx context.Context, name string, value string) error

	// GetValue 根据名称获取服务端配置值
	GetValue(ctx context.Context, name string) (string, error)

	// Exists 检查服务端配置是否存在
	Exists(ctx context.Context, name string) (bool, error)

	// BatchGetValues 批量获取服务端配置值
	BatchGetValues(ctx context.Context, names []string) (map[string]string, error)

	// Get 获取服务端配置，解析为指定类型
	Get(ctx context.Context, name string, value interface{}) error

	// Set 设置服务端配置，支持结构体类型
	Set(ctx context.Context, name string, value interface{}) error
}

// serverSettingServiceImpl 服务端设置服务实现
type serverSettingServiceImpl struct {
	repo serverRepo.ServerSettingRepository
}

// NewServerSettingService 创建服务端设置服务
func NewServerSettingService() ServerSettingService {
	return &serverSettingServiceImpl{
		repo: serverRepo.NewServerSettingRepository(db.GetDB()),
	}
}

// GetByName 根据名称获取服务端配置
func (s *serverSettingServiceImpl) GetByName(ctx context.Context, name string) (*server.ServerSetting, error) {
	return s.repo.GetByName(ctx, name)
}

// SetByName 设置服务端配置
func (s *serverSettingServiceImpl) SetByName(ctx context.Context, name string, value string) error {
	setting, err := s.GetByName(ctx, name)
	if err != nil {
		return err
	}

	if setting == nil {
		// 创建新配置
		setting = &server.ServerSetting{
			Name:  name,
			Value: value,
		}
		err = s.repo.Create(ctx, setting)
		return err
	}

	// 更新配置
	setting.Value = value
	return s.repo.Update(ctx, setting)
}

// GetValue 根据名称获取服务端配置值
func (s *serverSettingServiceImpl) GetValue(ctx context.Context, name string) (string, error) {
	setting, err := s.GetByName(ctx, name)
	if err != nil {
		return "", err
	}

	if setting == nil {
		return "", nil
	}

	return setting.Value, nil
}

// Exists 检查服务端配置是否存在
func (s *serverSettingServiceImpl) Exists(ctx context.Context, name string) (bool, error) {
	return s.repo.Exists(ctx, name)
}

// BatchGetValues 批量获取服务端配置值
func (s *serverSettingServiceImpl) BatchGetValues(ctx context.Context, names []string) (map[string]string, error) {
	settings, err := s.repo.BatchGet(ctx, names)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string, len(settings))
	for _, setting := range settings {
		result[setting.Name] = setting.Value
	}

	return result, nil
}

// Get 获取服务端配置，解析为指定类型
func (s *serverSettingServiceImpl) Get(ctx context.Context, name string, value interface{}) error {
	val, err := s.GetValue(ctx, name)
	if err != nil {
		return err
	}

	if val == "" {
		return nil
	}

	return json.Unmarshal([]byte(val), value)
}

// Set 设置服务端配置，支持结构体类型
func (s *serverSettingServiceImpl) Set(ctx context.Context, name string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.SetByName(ctx, name, string(data))
}
