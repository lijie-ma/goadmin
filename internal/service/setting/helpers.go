package setting

import (
	"encoding/json"
	appctx "goadmin/internal/context"
	"goadmin/internal/model/server"
)

// GetCaptchaSwitch 获取验证码开关状态
func GetCaptchaSwitch(ctx *appctx.Context, service ServerSettingService) (*server.CaptchaSwitchConfig, error) {
	value, err := service.GetValue(ctx.Request.Context(), server.SettingCaptchaSwitch)
	if err != nil {
		return nil, err
	}

	var config server.CaptchaSwitchConfig
	// 解析JSON
	err = json.Unmarshal([]byte(value), &config)
	if err != nil {
		return nil, err // 出错时返回默认配置
	}

	return &config, nil
}

// SetCaptchaSwitch 设置验证码开关状态
func SetCaptchaSwitch(ctx *appctx.Context, service ServerSettingService, config *server.CaptchaSwitchConfig) error {
	return service.Set(ctx.Request.Context(), server.SettingCaptchaSwitch, config)
}
