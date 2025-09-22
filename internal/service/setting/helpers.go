package setting

import (
	"goadmin/internal/context"
	"goadmin/internal/model/server"
)

// GetCaptchaSwitch 获取验证码开关状态
// 默认为开启状态
func GetCaptchaSwitch(ctx *context.Context, service ServerSettingService) (*server.CaptchaSwitch, error) {
	// 如果配置不存在，默认为开启
	var setting server.CaptchaSwitch
	err := service.Get(ctx, server.SettingCaptchaSwitch, &setting)
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// SetCaptchaSwitch 设置验证码开关状态
func SetCaptchaSwitch(ctx *context.Context, service ServerSettingService, setting *server.CaptchaSwitch) error {
	return service.Set(ctx, server.SettingCaptchaSwitch, setting)
}
