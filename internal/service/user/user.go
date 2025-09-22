package user

import (
	"fmt"
	"goadmin/internal/context"
	modeluser "goadmin/internal/model/user"
	userrepo "goadmin/internal/repository/user"
	"goadmin/internal/service/setting"
	"goadmin/internal/service/token"
	"goadmin/pkg/util"

	"goadmin/config"
)

// UserService 用户服务接口
type UserService interface {
	Login(ctx *context.Context, req modeluser.LoginRequest) (*modeluser.LoginResponse, error)
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

func (s *userService) Login(ctx *context.Context, req modeluser.LoginRequest) (*modeluser.LoginResponse, error) {
	captchaCfg, err := setting.GetCaptchaSwitch(ctx, setting.NewServerSettingService())
	if err != nil {
		ctx.Logger.Errorf("%s GetCaptchaSwitch %+v", s.logPrefix(), err)
		return nil, err
	}
	if captchaCfg.IsAdminOn() && !token.NewTokenService().ValidateToken(ctx, req.Token) {
		ctx.Logger.Errorf("%s ValidateToken faild %s %s", s.logPrefix(), req.Username, req.Token)
		return nil, err
	}
	// 获取用户信息
	u, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		ctx.Logger.Errorf("%s 获取用户信息失败: %s %v", s.logPrefix(), req.Username, err)
		return nil, err
	}

	if u == nil || !util.ValidatePasswordAndHash(req.Password, u.Password) {
		ctx.Logger.Warnf("%s 用户名或密码错误: %s %v", s.logPrefix(), req.Username, err)
		return nil, fmt.Errorf("用户名或密码错误")
	}

	if !u.IsActive() {
		ctx.Logger.Warnf("%s 账户状态异常: %s %s", s.logPrefix(), req.Username, u.Status.String())
		return nil, fmt.Errorf("账户状态异常")
	}

	// 生成JWT令牌
	tokenPairs, err := token.NewJwtTokenService(&s.cfg.JWT).GenerateJWTTokenPair(
		ctx, token.NewAdminClaims(u.ID, s.cfg.JWT.AccessExpire))
	if err != nil {
		ctx.Logger.Errorf("%s jwt token: %s %v", s.logPrefix(), req.Username, err)
		return nil, err
	}

	return &modeluser.LoginResponse{
		Token:        tokenPairs.AccessToken,
		RefreshToken: tokenPairs.RefreshToken,
		ExpiresAt:    tokenPairs.ExpiresAt,
		Role:         u.Role,
		Username:     u.Username,
		RoleCode:     u.RoleCode,
		Email:        u.Email,
	}, nil
}
