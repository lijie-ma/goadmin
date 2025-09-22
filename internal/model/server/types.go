package server

type Switch int8

const (
	// SwitchOff 关闭
	SwitchOff Switch = iota
	// SwitchOn 开启
	SwitchOn
)

// CaptchaSwitch 验证码开关
type CaptchaSwitch struct {
	Admin Switch `json:"admin"`
	Web   Switch `json:"web"` // web端
}

func (c CaptchaSwitch) IsAdminOn() bool {
	return c.Admin == SwitchOn
}

func (c CaptchaSwitch) IsWebOn() bool {
	return c.Web == SwitchOn
}
