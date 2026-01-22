package user

// UserStatus 用户状态
type UserStatus int

const (
	// UserStatusInactive 未激活
	UserStatusInactive UserStatus = iota
	// UserStatusActive 正常状态（从1开始）
	UserStatusActive // 值为1
	// UserStatusLocked 锁定状态
	UserStatusLocked // 值为2
	// UserStatusDeleted 已删除
	UserStatusDeleted // 值为3
)

// String 返回用户状态的字符串表示
func (s UserStatus) String() string {
	switch s {
	case UserStatusInactive:
		return "未激活"
	case UserStatusActive:
		return "正常"
	case UserStatusLocked:
		return "锁定"
	case UserStatusDeleted:
		return "已删除"
	default:
		return "未知状态"
	}
}

const (
	SupAdminUserID  = 1                                  // 超级管理员用户ID
	DefaultPassword = "e10adc3949ba59abbe56e057f20f883e" // 默认密码 123456
)
