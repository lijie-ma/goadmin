package server

type Switch int8

const (
	// SwitchOff 关闭
	SwitchOff Switch = iota
	// SwitchOn 开启
	SwitchOn
)

// CaptchaSwitch 验证码开关
type CaptchaSwitchConfig struct {
	Admin Switch `json:"admin" binding:"gte=0,lte=1"` // 管理后台开关
	Web   Switch `json:"web" binding:"gte=0,lte=1"`   // 网页端开关
}

// IsAdminOn 管理后台验证码是否开启
func (c *CaptchaSwitchConfig) IsAdminOn() bool {
	return c.Admin == SwitchOn
}

// IsWebOn 网页端验证码是否开启
func (c *CaptchaSwitchConfig) IsWebOn() bool {
	return c.Web == SwitchOn
}

type SystemConfig struct {
	SystemName string `json:"system_name" binding:"required"`
	Logo       string `json:"logo"`
	Language   string `json:"language" binding:"required"`
}

// CaptchaSwitchConfig 验证码开关配置

// SystemSettings 系统设置
type SystemSettingsRequest SystemSettingsResponse

type SystemSettingsResponse struct {
	SystemConfig
	CaptchaSwitchConfig
}
