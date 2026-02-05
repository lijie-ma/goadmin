package setting

import (
	"goadmin/internal/context"
	"goadmin/internal/i18n"
	"goadmin/internal/model/server"
	serverRepo "goadmin/internal/repository/server"
	"goadmin/pkg/db"
	"goadmin/pkg/util"
	"reflect"
)

// ServerSettingService 服务端设置服务接口
type ServerSettingService interface {
	// SetByName 设置服务端配置
	SetByName(ctx *context.Context, name string, value any) error

	// GetSrcValue 根据名称获取服务端配置值
	GetSrcValue(ctx *context.Context, name string, resultPtr any) error

	// GetValues 根据名称批量获取服务端配置值
	GetValues(ctx *context.Context, names []string) (map[string]any, error)

	// GetByName 根据名称获取服务端配置
	GetByName(ctx *context.Context, name string, decrypt ...bool) (*server.ServerSetting, error)

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

func (s *serverSettingServiceImpl) GetByName(ctx *context.Context, name string, decrypt ...bool) (*server.ServerSetting, error) {
	data, err := s.repo.GetByName(ctx, name)
	if err != nil {
		ctx.Logger.Errorf("%s getByName failed, err: %v", s.logPrefix(), err)
		return nil, i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}
	if data == nil {
		return nil, i18n.E(
			ctx.Context, "common.NotFound", map[string]any{"item": i18n.T(ctx.Context, "common.item.setting", nil)})
	}
	if len(decrypt) == 0 || decrypt[0] {
		decryptStr, err := util.DecryptAESGCM(data.Value)
		if err != nil {
			ctx.Logger.Errorf("%s SetByName DecryptAESGCM failed, err: %v", s.logPrefix(), err)
			return nil, err
		}
		data.Value = string(decryptStr)
	}
	return data, nil
}

// SetByName 设置服务端配置
func (s *serverSettingServiceImpl) SetByName(ctx *context.Context, name string, value any) error {
	str, err := encoding(value)
	if err != nil {
		return err
	}
	encryptStr, err := util.EncryptAESGCM([]byte(str))
	if err != nil {
		ctx.Logger.Errorf("%s SetByName EncryptAESGCM failed, err: %v", s.logPrefix(), err)
		return err
	}
	setting, err := s.repo.GetByName(ctx, name)
	if err != nil {
		ctx.Logger.Errorf("%s SetByName GetByName failed, err: %v", s.logPrefix(), err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}
	if setting == nil {
		// 创建新配置
		setting = &server.ServerSetting{
			Name:  name,
			Value: encryptStr,
		}
		err = s.repo.Create(ctx, setting)
	} else {
		setting.Value = encryptStr
		err = s.repo.Update(ctx, setting)
	}
	if err != nil {
		ctx.Logger.Errorf("%s SetByName failed, err: %v", s.logPrefix(), err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}

	return nil
}

// GetValue 根据名称获取服务端配置值
func (s *serverSettingServiceImpl) GetSrcValue(ctx *context.Context, name string, resultPtr any) error {
	setting, err := s.GetByName(ctx, name, false)
	if err != nil {
		return err
	}
	// 如果配置不存在，返回空对象
	if setting == nil {
		// 将resultPtr设置为空map
		if reflect.ValueOf(resultPtr).Kind() == reflect.Ptr {
			// 创建一个空的map[string]any并赋值给resultPtr
			emptyMap := make(map[string]any)
			reflect.ValueOf(resultPtr).Elem().Set(reflect.ValueOf(emptyMap))
		}
		return nil
	}
	return decoding(setting.Value, resultPtr)
}

// GetValues 根据名称批量获取服务端配置值
func (s *serverSettingServiceImpl) GetValues(ctx *context.Context, names []string) (map[string]any, error) {
	result := make(map[string]any)

	// 批量获取配置
	settings, err := s.repo.BatchGet(ctx, names)
	if err != nil {
		ctx.Logger.Errorf("%s GetValues failed, err: %v", s.logPrefix(), err)
		return nil, i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}

	// 将配置转换为map
	for _, setting := range settings {
		var value any
		err = decoding(setting.Value, &value)
		if err != nil {
			ctx.Logger.Errorf("%s GetValues unmarshal failed for %s, err: %v", s.logPrefix(), setting.Name, err)
			// 如果解析失败，返回空对象
			result[setting.Name] = make(map[string]any)
			continue
		}
		result[setting.Name] = value
	}

	// 对于不存在的配置，返回空对象
	for _, name := range names {
		if _, exists := result[name]; !exists {
			result[name] = make(map[string]any)
		}
	}

	return result, nil
}

// GetSystemSettings 获取系统设置
func (s *serverSettingServiceImpl) GetSystemSettings(ctx *context.Context) (*server.SystemSettingsResponse, error) {
	cfgs, err := s.repo.BatchGet(ctx, []string{server.SettingCaptchaSwitch, server.SettingSystemConfig})
	if err != nil {
		ctx.Logger.Errorf("%s GetSystemSettings failed, err: %v", s.logPrefix(), err)
		return nil, i18n.E(ctx.Context, "common.RepositoryErr", nil)
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
		return err
	}
	err = s.SetByName(ctx, server.SettingSystemConfig, settings.SystemConfig)
	if err != nil {
		ctx.Logger.Errorf("%s SetSystemSettings SettingSystemConfig failed, err: %v", s.logPrefix(), err)
		return err
	}
	return nil
}
