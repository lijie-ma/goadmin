package user

import (
	"goadmin/internal/context"
	"goadmin/internal/i18n"
	"goadmin/internal/model/schema"
	"goadmin/internal/model/server"
	modeluser "goadmin/internal/model/user"
	"goadmin/internal/repository/role"
	userrepo "goadmin/internal/repository/user"
	"goadmin/internal/service/setting"
	"goadmin/internal/service/token"
	"goadmin/pkg/util"

	"goadmin/config"
)

// UserService 用户服务接口
type UserService interface {
	// Login 用户登录
	Login(ctx *context.Context, req modeluser.LoginRequest) (*modeluser.LoginResponse, error)

	// GenerateUserCredential 根据用户ID生成用户身份凭证
	GenerateUserCredential(ctx *context.Context, userID uint64) (*token.TokenPair, error)

	// GenerateUserSession 当前Session
	GenerateUserSession(ctx *context.Context, userID uint64) (*modeluser.User, error)

	// ListUsers 获取用户列表
	ListUsers(ctx *context.Context, req *modeluser.ListRequest) ([]*modeluser.User, int64, error)

	// CreateUser 创建用户
	CreateUser(ctx *context.Context, req *modeluser.CreateUserRequest) error

	// UpdateUser 更新用户
	UpdateUser(ctx *context.Context, req *modeluser.UpdateUserRequest) error

	// DeleteUser 删除用户
	DeleteUser(ctx *context.Context, req *schema.IDRequest) error

	ChangePassword(ctx *context.Context, req *modeluser.ChangePasswordRequest) error

	// ResetPassword 重置密码
	ResetPassword(ctx *context.Context, req *schema.IDRequest) error
}

// userService 用户服务实现
type userService struct {
	userRepo userrepo.UserRepository
	roleRepo role.RoleRepository
	cfg      *config.Config
}

// NewUserService 创建用户服务实例
func NewUserService() UserService {
	return &userService{
		cfg:      config.Get(),
		userRepo: userrepo.NewUserRepository(),
		roleRepo: role.NewRoleRepositoryWithDB(),
	}
}

func (*userService) logPrefix() string {
	return "user-service"
}

// GenerateUserCredential 根据用户ID生成用户身份凭证
func (s *userService) GenerateUserCredential(ctx *context.Context, userID uint64) (*token.TokenPair, error) {
	// 生成JWT令牌
	tokenPairs, err := token.NewJwtTokenService(&s.cfg.JWT).GenerateJWTTokenPair(
		ctx, token.NewAdminClaims(userID, s.cfg.JWT.AccessExpire))
	if err != nil {
		ctx.Logger.Errorf("%s 生成用户凭证失败: %d %v", s.logPrefix(), userID, err)
		return nil, i18n.E(ctx.Context, "user.token.generate.failed", nil)
	}

	return tokenPairs, nil
}

func (s *userService) Login(ctx *context.Context, req modeluser.LoginRequest) (*modeluser.LoginResponse, error) {
	var captchaCfg server.CaptchaSwitchConfig
	err := setting.NewServerSettingService().GetValue(ctx, server.SettingCaptchaSwitch, &captchaCfg)
	if err != nil {
		ctx.Logger.Errorf("%s Generate GetValue %+v", s.logPrefix(), err)
		return nil, i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}
	if captchaCfg.IsAdminOn() {
		if req.Token == "" {
			ctx.Logger.Warnf("%s Login captcha require token", s.logPrefix())
			return nil, i18n.E(ctx.Context, "common.BadParameter", nil)
		}
		if !token.NewTokenService().ValidateToken(ctx, req.Token) {
			ctx.Logger.Errorf("%s ValidateToken faild %s %s", s.logPrefix(), req.Username, req.Token)
			return nil, i18n.E(
				ctx.Context,
				"common.InvalidParameter",
				map[string]any{"item": i18n.T(ctx.Context, "common.item.user", nil)})
		}
	}
	// 获取用户信息
	u, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		ctx.Logger.Errorf("%s 获取用户信息失败: %s %v", s.logPrefix(), req.Username, err)
		return nil, err
	}

	if u == nil || !util.ValidatePasswordAndHash(req.Password, u.Password) {
		ctx.Logger.Warnf("%s 用户名或密码错误: %s %v", s.logPrefix(), req.Username, err)
		return nil, i18n.E(ctx.Context, "user.InvalidUsernameOrPassword", nil)
	}

	if !u.IsActive() {
		ctx.Logger.Warnf("%s 账户状态异常: %s %s", s.logPrefix(), req.Username, u.Status.String())
		return nil, i18n.E(ctx.Context, "user.AccountStatusAbnormal", nil)
	}

	// 生成JWT令牌
	tokenPairs, err := s.GenerateUserCredential(ctx, u.ID)
	if err != nil {
		ctx.Logger.Errorf("%s jwt token: %s %v", s.logPrefix(), req.Username, err)
		return nil, err
	}
	perms, err := role.NewRolePermissionRepositoryWithDB().GetPermissionsByRoleCode(ctx, u.RoleCode)
	if err != nil {
		ctx.Logger.Errorf("%s get Permissions: %s %d %v", s.logPrefix(), req.Username, u.RoleCode, err)
		return nil, i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}

	return &modeluser.LoginResponse{
		Token:           tokenPairs.AccessToken,
		RefreshToken:    tokenPairs.RefreshToken,
		ExpiresAt:       tokenPairs.ExpiresAt,
		Role:            u.Role,
		Username:        u.Username,
		RoleCode:        u.RoleCode,
		Email:           u.Email,
		PermissionCodes: perms,
	}, nil
}

// GenerateUserSession 当前Session
func (s *userService) GenerateUserSession(ctx *context.Context, userID uint64) (*modeluser.User, error) {
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		ctx.Logger.Errorf("%s 获取用户信息失败: %d %v", s.logPrefix(), userID, err)
		return nil, i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}

	if u == nil {
		ctx.Logger.Warnf("%s 用户不存在: %d", s.logPrefix(), userID)
		return nil, i18n.E(ctx.Context, "common.NotFound", map[string]any{"item": i18n.T(ctx.Context, "common.item.user", nil)})
	}

	if !u.IsActive() {
		ctx.Logger.Warnf("%s 账户状态异常: %s %s", s.logPrefix(), userID, u.Status.String())
		return nil, i18n.E(ctx.Context, "user.AccountStatusAbnormal", nil)
	}
	return u, nil
}

// ListUsers 获取用户列表
func (s *userService) ListUsers(ctx *context.Context, req *modeluser.ListRequest) ([]*modeluser.User, int64, error) {
	list, total, err := s.userRepo.PageList(ctx, req)
	if err != nil {
		ctx.Logger.Errorf("%s 获取用户列表失败: %v", s.logPrefix(), err)
		return nil, 0, i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}
	if total == 0 {
		return []*modeluser.User{}, 0, nil
	}
	return list, total, nil
}

// CreateUser 创建用户
func (s *userService) CreateUser(ctx *context.Context, req *modeluser.CreateUserRequest) error {
	// 检查用户名是否已存在
	exists, err := s.userRepo.IsUsernameExists(ctx, req.Username)
	if err != nil {
		ctx.Logger.Errorf("%s 检查用户名是否存在失败: %s %v", s.logPrefix(), req.Username, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}
	if exists {
		ctx.Logger.Warnf("%s 用户名已存在: %s", s.logPrefix(), req.Username)
		return i18n.E(ctx.Context, "user.UsernameAlreadyExists", nil)
	}

	// 检查邮箱是否已存在
	if req.Email != "" {
		exists, err = s.userRepo.IsEmailExists(ctx, req.Email)
		if err != nil {
			ctx.Logger.Errorf("%s 检查邮箱是否存在失败: %s %v", s.logPrefix(), req.Email, err)
			return i18n.E(ctx.Context, "common.RepositoryErr", nil)
		}
		if exists {
			ctx.Logger.Warnf("%s 邮箱已存在: %s", s.logPrefix(), req.Email)
			return i18n.E(ctx.Context, "user.EmailAlreadyExists", nil)
		}
	}

	// 检查角色是否存在
	role, err := s.roleRepo.GetByCode(ctx, req.RoleCode)
	if err != nil {
		ctx.Logger.Errorf("%s 检查角色是否存在失败: %s %v", s.logPrefix(), req.RoleCode, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}
	if role == nil {
		ctx.Logger.Warnf("%s 角色不存在: %s", s.logPrefix(), req.RoleCode)
		return i18n.E(ctx.Context, "common.NotFound", map[string]any{"item": i18n.T(ctx.Context, "common.item.role", nil)})
	}

	// 加密密码
	encryptPwd, err := util.Password2Hash(req.Password)
	if err != nil {
		ctx.Logger.Errorf("%s 密码加密失败: %s %v", s.logPrefix(), req.Username, err)
		return i18n.E(ctx.Context, "common.EncryptErr", nil)
	}

	// 创建用户
	user := &modeluser.User{
		Username: req.Username,
		Password: encryptPwd,
		Email:    req.Email,
		RoleCode: req.RoleCode,
		Status:   modeluser.UserStatus(req.Status),
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		ctx.Logger.Errorf("%s 创建用户失败: %s %v", s.logPrefix(), req.Username, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}

	ctx.Logger.Infof("%s 创建用户成功: %s", s.logPrefix(), req.Username)
	return nil
}

// UpdateUser 更新用户
func (s *userService) UpdateUser(ctx *context.Context, req *modeluser.UpdateUserRequest) error {
	// 获取用户信息
	user, err := s.userRepo.GetByID(ctx, req.ID)
	if err != nil {
		ctx.Logger.Errorf("%s 获取用户信息失败: %d %v", s.logPrefix(), req.ID, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}
	if user == nil {
		ctx.Logger.Warnf("%s 用户不存在: %d", s.logPrefix(), req.ID)
		return i18n.E(ctx.Context, "common.NotFound", map[string]any{"item": i18n.T(ctx.Context, "common.item.user", nil)})
	}

	// 检查是否为超级管理员
	if user.IsSuperAdmin() {
		ctx.Logger.Warnf("%s 超级管理员不能修改: %d", s.logPrefix(), req.ID)
		return i18n.E(ctx.Context, "user.SuperAdminCannotModify", nil)
	}

	// 检查用户名是否已被其他用户使用
	if req.Username != "" && req.Username != user.Username {
		exists, err := s.userRepo.IsUsernameExists(ctx, req.Username, req.ID)
		if err != nil {
			ctx.Logger.Errorf("%s 检查用户名是否存在失败: %s %v", s.logPrefix(), req.Username, err)
			return i18n.E(ctx.Context, "common.RepositoryErr", nil)
		}
		if exists {
			ctx.Logger.Warnf("%s 用户名已存在: %s", s.logPrefix(), req.Username)
			return i18n.E(ctx.Context, "user.UsernameAlreadyExists", nil)
		}
		user.Username = req.Username
	}

	// 检查邮箱是否已被其他用户使用
	if req.Email != "" && req.Email != user.Email {
		exists, err := s.userRepo.IsEmailExists(ctx, req.Email, req.ID)
		if err != nil {
			ctx.Logger.Errorf("%s 检查邮箱是否存在失败: %s %v", s.logPrefix(), req.Email, err)
			return i18n.E(ctx.Context, "common.RepositoryErr", nil)
		}
		if exists {
			ctx.Logger.Warnf("%s 邮箱已存在: %s", s.logPrefix(), req.Email)
			return i18n.E(ctx.Context, "user.EmailAlreadyExists", nil)
		}
		user.Email = req.Email
	}

	// 更新角色代码
	if req.RoleCode != "" && req.RoleCode != user.RoleCode {
		// 检查角色是否存在
		role, err := s.roleRepo.GetByCode(ctx, req.RoleCode)
		if err != nil {
			ctx.Logger.Errorf("%s 检查角色是否存在失败: %s %v", s.logPrefix(), req.RoleCode, err)
			return i18n.E(ctx.Context, "common.RepositoryErr", nil)
		}
		if role == nil {
			ctx.Logger.Warnf("%s 角色不存在: %s", s.logPrefix(), req.RoleCode)
			return i18n.E(ctx.Context, "common.NotFound", map[string]any{"item": i18n.T(ctx.Context, "common.item.role", nil)})
		}
		user.RoleCode = req.RoleCode
	}

	// 更新状态
	if req.Status >= 0 {
		user.Status = modeluser.UserStatus(req.Status)
	}

	// 更新用户
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		ctx.Logger.Errorf("%s 更新用户失败: %d %v", s.logPrefix(), req.ID, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}

	ctx.Logger.Infof("%s 更新用户成功: %d", s.logPrefix(), req.ID)
	return nil
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(ctx *context.Context, req *schema.IDRequest) error {
	// 获取用户信息
	user, err := s.userRepo.GetByID(ctx, req.ID)
	if err != nil {
		ctx.Logger.Errorf("%s 获取用户信息失败: %d %v", s.logPrefix(), req.ID, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}
	if user == nil {
		ctx.Logger.Warnf("%s 用户不存在: %d", s.logPrefix(), req.ID)
		return i18n.E(
			ctx.Context, "common.NotFound",
			map[string]any{"item": i18n.T(ctx.Context, "common.item.user", nil)})
	}

	// 检查是否为超级管理员
	if user.IsSuperAdmin() {
		ctx.Logger.Warnf("%s 超级管理员不能删除: %d", s.logPrefix(), req.ID)
		return i18n.E(ctx.Context, "user.SuperAdminCannotDelete", nil)
	}

	// 检查是否为当前登录用户
	if req.ID == ctx.Session().GetID() {
		ctx.Logger.Warnf("%s 不能删除当前登录用户: %d", s.logPrefix(), req.ID)
		return i18n.E(ctx.Context, "user.CannotDeleteSelf", nil)
	}

	// 软删除用户（将状态设置为已删除）
	err = s.userRepo.UpdateStatus(ctx, req.ID, modeluser.UserStatusDeleted)
	if err != nil {
		ctx.Logger.Errorf("%s 删除用户失败: %d %v", s.logPrefix(), req.ID, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}

	ctx.Logger.Infof("%s 删除用户成功: %d", s.logPrefix(), req.ID)
	return nil
}

func (s *userService) ChangePassword(ctx *context.Context, req *modeluser.ChangePasswordRequest) error {
	if !util.ValidatePasswordAndHash(req.OldPassword, ctx.Session().(*modeluser.User).Password) {
		ctx.Logger.Warnf("%s 密码错误: %s", s.logPrefix(), ctx.Session().GetUsername())
		return i18n.E(ctx.Context, "user.InvalidPassword", nil)
	}

	encryptPwd, err := util.Password2Hash(req.NewPassword)
	if err != nil {
		ctx.Logger.Warnf("%s Password2Hash: %s %+v", s.logPrefix(), ctx.Session().GetUsername(), err)
		return i18n.E(ctx.Context, "common.EncryptErr", nil)
	}

	// 更新密码
	err = s.userRepo.UpdatePassword(ctx, ctx.Session().GetID(), encryptPwd)
	if err != nil {
		ctx.Logger.Warnf("%s UpdatePassword: %s %+v", s.logPrefix(), ctx.Session().GetUsername(), err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}
	return nil
}

// ResetPassword 重置密码
func (s *userService) ResetPassword(ctx *context.Context, req *schema.IDRequest) error {
	// 获取用户信息
	user, err := s.userRepo.GetByID(ctx, req.ID)
	if err != nil {
		ctx.Logger.Errorf("%s 获取用户信息失败: %d %v", s.logPrefix(), req.ID, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}
	if user == nil {
		ctx.Logger.Warnf("%s 用户不存在: %d", s.logPrefix(), req.ID)
		return i18n.E(ctx.Context, "common.NotFound", map[string]any{"item": i18n.T(ctx.Context, "common.item.user", nil)})
	}

	// 检查是否为超级管理员
	if user.IsSuperAdmin() {
		ctx.Logger.Warnf("%s 超级管理员不能重置密码: %d", s.logPrefix(), req.ID)
		return i18n.E(ctx.Context, "user.SuperAdminCannotModify", nil)
	}

	// 加密默认密码
	encryptPwd, err := util.Password2Hash(modeluser.DefaultPassword)
	if err != nil {
		ctx.Logger.Errorf("%s 密码加密失败: %d %v", s.logPrefix(), req.ID, err)
		return i18n.E(ctx.Context, "common.EncryptErr", nil)
	}

	// 更新密码
	err = s.userRepo.UpdatePassword(ctx, req.ID, encryptPwd)
	if err != nil {
		ctx.Logger.Errorf("%s 重置密码失败: %d %v", s.logPrefix(), req.ID, err)
		return i18n.E(ctx.Context, "common.RepositoryErr", nil)
	}

	ctx.Logger.Infof("%s 重置密码成功: %d", s.logPrefix(), req.ID)
	return nil
}
