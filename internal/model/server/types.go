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
	Admin  Switch `json:"admin"`  // 管理后台开关
	Mobile Switch `json:"mobile"` // 移动端开关
	Web    Switch `json:"web"`    // 网页端开关
}

// IsAdminOn 管理后台验证码是否开启
func (c *CaptchaSwitchConfig) IsAdminOn() bool {
	return c.Admin == SwitchOn
}

// IsMobileOn 移动端验证码是否开启
func (c *CaptchaSwitchConfig) IsMobileOn() bool {
	return c.Mobile == SwitchOn
}

// IsWebOn 网页端验证码是否开启
func (c *CaptchaSwitchConfig) IsWebOn() bool {
	return c.Web == SwitchOn
}
