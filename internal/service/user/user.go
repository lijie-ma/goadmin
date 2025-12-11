package user

import (
	"fmt"
	"goadmin/internal/context"
	"goadmin/internal/model/server"
	modeluser "goadmin/internal/model/user"
	"goadmin/internal/repository/role"
	userrepo "goadmin/internal/repository/user"
	"goadmin/internal/service/errorsx"
	"goadmin/internal/service/setting"
	"goadmin/internal/service/token"
	"goadmin/pkg/util"

	"goadmin/config"

	"github.com/pkg/errors"
)

// UserService 用户服务接口
type UserService interface {
	// Login 用户登录
	Login(ctx *context.Context, req modeluser.LoginRequest) (*modeluser.LoginResponse, error)

	// GenerateUserCredential 根据用户ID生成用户身份凭证
	GenerateUserCredential(ctx *context.Context, userID uint64) (*token.TokenPair, error)

	// GenerateUserSession 当前Session
	GenerateUserSession(ctx *context.Context, userID uint64) (*modeluser.User, error)

	ChangePassword(ctx *context.Context, req *modeluser.ChangePasswordRequest) error
}

// userService 用户服务实现
type userService struct {
	userRepo userrepo.UserRepository
	cfg      *config.Config
}

// NewUserService 创建用户服务实例
func NewUserService() UserService {
	return &userService{
		cfg:      config.Get(),
		userRepo: userrepo.NewUserRepository(),
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
		return nil, err
	}

	return tokenPairs, nil
}

func (s *userService) Login(ctx *context.Context, req modeluser.LoginRequest) (*modeluser.LoginResponse, error) {
	var captchaCfg server.CaptchaSwitchConfig
	err := setting.NewServerSettingService().GetValue(ctx, server.SettingCaptchaSwitch, &captchaCfg)
	if err != nil {
		ctx.Logger.Errorf("%s Generate GetValue %+v", s.logPrefix(), err)
		return nil, err
	}
	if captchaCfg.IsAdminOn() {
		if req.Token == "" {
			return nil, errors.WithMessage(errorsx.ErrReqired, "token")
		}
		if !token.NewTokenService().ValidateToken(ctx, req.Token) {
			ctx.Logger.Errorf("%s ValidateToken faild %s %s", s.logPrefix(), req.Username, req.Token)
			return nil, errors.WithMessage(errorsx.ErrInvalid, "token")
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
		return nil, errors.WithMessage(errorsx.ErrNotFound, "用户")
	}

	if !u.IsActive() {
		ctx.Logger.Warnf("%s 账户状态异常: %s %s", s.logPrefix(), req.Username, u.Status.String())
		return nil, fmt.Errorf("账户状态异常")
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
		return nil, err
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
		return nil, err
	}

	if !u.IsActive() {
		ctx.Logger.Warnf("%s 账户状态异常: %s %s", s.logPrefix(), userID, u.Status.String())
		return nil, fmt.Errorf("账户状态异常")
	}
	return u, nil
}

func (s *userService) ChangePassword(ctx *context.Context, req *modeluser.ChangePasswordRequest) error {
	if !util.ValidatePasswordAndHash(req.OldPassword, ctx.Session().(*modeluser.User).Password) {
		ctx.Logger.Warnf("%s 密码错误: %s", s.logPrefix(), ctx.Session().GetUsername())
		return fmt.Errorf("密码错误")
	}

	encryptPwd, err := util.Password2Hash(req.NewPassword)
	if err != nil {
		ctx.Logger.Warnf("%s Password2Hash: %s %+v", s.logPrefix(), ctx.Session().GetUsername(), err)
		return err
	}

	// 更新密码
	err = s.userRepo.UpdatePassword(ctx, ctx.Session().GetID(), encryptPwd)
	if err != nil {
		ctx.Logger.Warnf("%s UpdatePassword: %s %+v", s.logPrefix(), ctx.Session().GetUsername(), err)
		return err
	}
	return nil
}
