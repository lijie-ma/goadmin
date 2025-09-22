package setting

import (
	"encoding/json"
	"goadmin/internal/context"
	"goadmin/internal/model/server"
	serverrepo "goadmin/internal/repository/server"
	"goadmin/pkg/db"
)

// ServerSettingService 服务端配置服务接口
type ServerSettingService interface {
	// Get 获取配置值
	//  值将存放到valuePointer 中
	Get(ctx *context.Context, name string, valuePointer any) error

	// Set 设置配置值
	//  value 将被 json 序列化
	Set(ctx *context.Context, name string, value any) error

	// GetBatch 批量获取配置值
	GetBatch(ctx *context.Context, names []string) (map[string]string, error)
}

// serverSettingService 服务端配置服务实现
type serverSettingService struct {
	settingRepo serverrepo.ServerSettingRepository
}

// NewServerSettingService 创建服务端配置服务实例
func NewServerSettingService() ServerSettingService {
	return &serverSettingService{
		settingRepo: serverrepo.NewServerSettingRepository(db.GetDB()), // 这里的 nil 会使用全局 DB 实例
	}
}

func (*serverSettingService) logPrefix() string {
	return "server-setting-service"
}

// Get 获取配置值
func (s *serverSettingService) Get(ctx *context.Context, name string, valuePointer any) error {
	setting, err := s.settingRepo.GetByName(ctx, name)
	if err != nil {
		ctx.Logger.Errorf("%s 获取配置失败: %s %v", s.logPrefix(), name, err)
		return err
	}

	if setting == nil {
		return nil
	}

	err = json.Unmarshal([]byte(setting.Value), valuePointer)
	if err != nil {
		ctx.Logger.Errorf("%s 解析配置值失败: %s %v", s.logPrefix(), name, err)
		return err
	}

	return nil
}

// Set 设置配置值
func (s *serverSettingService) Set(ctx *context.Context, name string, value any) error {
	str, err := json.Marshal(value)
	if err != nil {
		ctx.Logger.Errorf("%s 序列化配置值失败: %s %v", s.logPrefix(), name, err)
		return err
	}
	// 先检查配置是否存在
	exists, err := s.settingRepo.Exists(ctx, name)
	if err != nil {
		ctx.Logger.Errorf("%s 检查配置是否存在失败: %s %v", s.logPrefix(), name, err)
		return err
	}

	setting := &server.ServerSetting{
		Name:  name,
		Value: string(str),
	}

	if exists {
		// 如果存在则更新
		err = s.settingRepo.Update(ctx, setting)
	} else {
		// 如果不存在则创建
		err = s.settingRepo.Create(ctx, setting)
	}

	if err != nil {
		ctx.Logger.Errorf("%s 保存配置失败: %s %v", s.logPrefix(), name, err)
		return err
	}

	return nil
}

// GetBatch 批量获取配置值
func (s *serverSettingService) GetBatch(ctx *context.Context, names []string) (map[string]string, error) {
	settings, err := s.settingRepo.BatchGet(ctx, names)
	if err != nil {
		ctx.Logger.Errorf("%s 批量获取配置失败: %v %v", s.logPrefix(), names, err)
		return nil, err
	}

	result := make(map[string]string, len(settings))
	for _, setting := range settings {
		result[setting.Name] = setting.Value
	}

	return result, nil
}
