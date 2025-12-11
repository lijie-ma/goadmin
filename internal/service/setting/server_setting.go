package setting

import (
	"fmt"
	"goadmin/internal/context"
	"goadmin/internal/model/server"
	serverRepo "goadmin/internal/repository/server"
	"goadmin/internal/service/errorsx"
	"goadmin/pkg/db"
)

// ServerSettingService 服务端设置服务接口
type ServerSettingService interface {
	// SetByName 设置服务端配置
	SetByName(ctx *context.Context, name string, value any) error

	// GetValue 根据名称获取服务端配置值
	GetValue(ctx *context.Context, name string, resultPtr any) error

	// GetByName 根据名称获取服务端配置
	GetByName(ctx *context.Context, name string) (*server.ServerSetting, error)

	// Exists 检查服务端配置是否存在
	Exists(ctx *context.Context, name string) (bool, error)

	// GetSystemSettings 获取系统设置
	GetSystemSettings(ctx *context.Context) (*server.SystemSettingsResponse, error)

	// SetSystemSettings 设置系统设置
	SetSystemSettings(ctx *context.Context, settings *server.SystemSettingsRequest) error
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

func (s *serverSettingServiceImpl) logPrefix() string {
	return "server-setting"
}

func (s *serverSettingServiceImpl) getByName(ctx *context.Context, name string) (*server.ServerSetting, error) {
	data, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, errorsx.ErrNotFound
	}
	return data, nil
}

// SetByName 设置服务端配置
func (s *serverSettingServiceImpl) SetByName(ctx *context.Context, name string, value any) error {
	setting, err := s.getByName(ctx, name)
	if err != nil && err != errorsx.ErrNotFound {
		return err
	}
	str, err := encoding(value)
	if err != nil {
		return err
	}
	if setting == nil {
		// 创建新配置
		setting = &server.ServerSetting{
			Name:  name,
			Value: string(str),
		}
		err = s.repo.Create(ctx, setting)
		return err
	}

	// 更新配置
	setting.Value = string(str)
	return s.repo.Update(ctx, setting)
}

// GetValue 根据名称获取服务端配置值
func (s *serverSettingServiceImpl) GetValue(ctx *context.Context, name string, resultPtr any) error {
	setting, err := s.getByName(ctx, name)
	if err != nil {
		return err
	}
	return decoding(setting.Value, resultPtr)
}

// GetByName 根据名称获取服务端配置
func (s *serverSettingServiceImpl) GetByName(ctx *context.Context, name string) (*server.ServerSetting, error) {
	return s.getByName(ctx, name)
}

// Exists 检查服务端配置是否存在
func (s *serverSettingServiceImpl) Exists(ctx *context.Context, name string) (bool, error) {
	return s.repo.Exists(ctx, name)
}

// GetSystemSettings 获取系统设置
func (s *serverSettingServiceImpl) GetSystemSettings(ctx *context.Context) (*server.SystemSettingsResponse, error) {
	cfgs, err := s.repo.BatchGet(ctx, []string{server.SettingCaptchaSwitch, server.SettingSystemConfig})
	if err != nil {
		ctx.Logger.Errorf("%s GetSystemSettings failed, err: %v", s.logPrefix(), err)
		return nil, err
	}
	var (
		rs            server.SystemSettingsResponse
		captchaSwitch server.CaptchaSwitchConfig
		systemConfig  server.SystemConfig
	)

	for _, cfg := range cfgs {
		switch cfg.Name {
		case server.SettingCaptchaSwitch:
			err = decoding(cfg.Value, &captchaSwitch)
			if err != nil {
				ctx.Logger.Errorf("%s GetSystemSettings unmarshal captcha switch failed, err: %v", s.logPrefix(), err)
				return nil, err
			}
		case server.SettingSystemConfig:
			err = decoding(cfg.Value, &systemConfig)
			if err != nil {
				ctx.Logger.Errorf("%s GetSystemSettings unmarshal system config failed, err: %v", s.logPrefix(), err)
				return nil, err
			}
		}
	}
	rs.CaptchaSwitchConfig = captchaSwitch
	rs.SystemConfig = systemConfig
	return &rs, nil
}

// SetSystemSettings 设置系统设置
func (s *serverSettingServiceImpl) SetSystemSettings(ctx *context.Context, settings *server.SystemSettingsRequest) error {
	err := s.SetByName(ctx, server.SettingCaptchaSwitch, settings.CaptchaSwitchConfig)
	if err != nil {
		ctx.Logger.Errorf("%s SetSystemSettings SetCaptchaSwitch failed, err: %v", s.logPrefix(), err)
		return fmt.Errorf("设置验证码开关失败: %w", err)
	}
	err = s.SetByName(ctx, server.SettingSystemConfig, settings.SystemConfig)
	if err != nil {
		ctx.Logger.Errorf("%s SetSystemSettings SettingSystemConfig failed, err: %v", s.logPrefix(), err)
		return fmt.Errorf("系统配置: %w", err)
	}
	return nil
}
