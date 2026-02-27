// Package wire provides dependency injection configuration using Google Wire.
package wire

import (
	"goadmin/config"

	// Infrastructure
	"goadmin/pkg/db"
	"goadmin/pkg/redisx"
	"goadmin/pkg/task"

	// Internal
	"goadmin/internal/api"
	"goadmin/internal/i18n"

	// Repository
	userrepo "goadmin/internal/repository/user"
	rolerepo "goadmin/internal/repository/role"
	operatelogrepo "goadmin/internal/repository/operate_log"
	positionrepo "goadmin/internal/repository/position"
	serverrepo "goadmin/internal/repository/server"

	// Service
	"goadmin/internal/service/token"
	"goadmin/internal/service/captcha"
	"goadmin/internal/service/operate_log"
	"goadmin/internal/service/position"
	"goadmin/internal/service/role"
	"goadmin/internal/service/setting"
	userservice "goadmin/internal/service/user"

	// HTTP Server
	serverpkg "goadmin/cmd/server"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/google/wire"
)

// ============================================================================
// Infrastructure Providers
// ============================================================================

// ProvideConfig provides the global configuration instance.
func ProvideConfig() *config.Config {
	return config.Get()
}

// ProvideDB provides the database connection.
func ProvideDB(cfg *config.Config) (*gorm.DB, error) {
	err := db.Init(&cfg.Database)
	if err != nil {
		return nil, err
	}
	return db.GetDB(), nil
}

// ProvideRedis provides the Redis client initialization.
func ProvideRedis(cfg *config.Config) error {
	return redisx.Init(&cfg.Redis)
}

// ============================================================================
// Repository Providers
// ============================================================================

// ProvideUserRepository provides the user repository.
func ProvideUserRepository(database *gorm.DB) userrepo.UserRepository {
	return userrepo.NewUserRepositoryImpl(database)
}

// ProvideRoleRepository provides the role repository.
func ProvideRoleRepository(database *gorm.DB) rolerepo.RoleRepository {
	return rolerepo.NewRoleRepository(database)
}

// ProvideRolePermissionRepository provides the role permission repository.
func ProvideRolePermissionRepository(database *gorm.DB) rolerepo.RolePermissionRepository {
	return rolerepo.NewRolePermissionRepository(database)
}

// ProvideOperateLogRepository provides the operate log repository.
func ProvideOperateLogRepository(database *gorm.DB) operatelogrepo.OperateLogRepository {
	return operatelogrepo.NewOperateLogRepositoryImpl(database)
}

// ProvidePositionRepository provides the position repository.
func ProvidePositionRepository(database *gorm.DB) positionrepo.PositionRepository {
	return positionrepo.NewPositionRepositoryImpl(database)
}

// ProvideServerSettingRepository provides the server setting repository.
func ProvideServerSettingRepository(database *gorm.DB) serverrepo.ServerSettingRepository {
	return serverrepo.NewServerSettingRepository(database)
}

// ============================================================================
// Service Providers
// ============================================================================

// ProvideTokenService provides the token service.
func ProvideTokenService() *token.TokenService {
	return token.NewTokenService()
}

// ProvideJwtTokenService provides the JWT token service.
func ProvideJwtTokenService(cfg *config.Config) *token.JwtTokenService {
	return token.NewJwtTokenService(&cfg.JWT)
}

// ProvideCaptchaService provides the captcha service.
func ProvideCaptchaService() captcha.CaptchaService {
	return captcha.NewCaptchaService()
}

// ProvideServerSettingService provides the server setting service.
func ProvideServerSettingService(repo serverrepo.ServerSettingRepository) setting.ServerSettingService {
	return setting.NewServerSettingService(repo)
}

// ProvideOperateLogService provides the operate log service.
func ProvideOperateLogService(logRepo operatelogrepo.OperateLogRepository) operate_log.OperateLogService {
	return operate_log.NewOperateLogService(logRepo)
}

// ProvidePositionService provides the position service.
func ProvidePositionService(positionRepo positionrepo.PositionRepository, logService operate_log.OperateLogService) position.PositionService {
	return position.NewPositionService(positionRepo, logService)
}

// ProvideRoleService provides the role service.
func ProvideRoleService(roleRepo rolerepo.RoleRepository, rolePermissionRepo rolerepo.RolePermissionRepository, cfg *config.Config) role.RoleService {
	return role.NewRoleService(roleRepo, rolePermissionRepo, cfg)
}

// ProvideUserService provides the user service.
func ProvideUserService(
	cfg *config.Config,
	userRepo userrepo.UserRepository,
	roleRepo rolerepo.RoleRepository,
	logService operate_log.OperateLogService,
	tokenService *token.TokenService,
	jwtTokenService *token.JwtTokenService,
	captchaService captcha.CaptchaService,
	serverSettingService setting.ServerSettingService,
) userservice.UserService {
	return userservice.NewUserService(
		cfg,
		userRepo,
		roleRepo,
		logService,
		tokenService,
		jwtTokenService,
		captchaService,
		serverSettingService,
	)
}

// ============================================================================
// HTTP Server Providers
// ============================================================================

// ProvideGinEngine provides the Gin HTTP engine.
// Note: Gin mode is set by NewWebServer
func ProvideGinEngine() *gin.Engine {
	// Initialize i18n to ensure Bundle is not nil
	i18n.Init()
	// Create Gin instance without default middleware
	// Middleware will be registered by RegisterRouter
	r := gin.New()
	return r
}

// ProvideWebServer provides the web server.
func ProvideWebServer(
	cfg *config.Config,
	engine *gin.Engine,
	tokenService *token.TokenService,
	userService userservice.UserService,
	roleService role.RoleService,
	positionService position.PositionService,
	logService operate_log.OperateLogService,
	settingService setting.ServerSettingService,
	userRepository userrepo.UserRepository,
) *serverpkg.WebServer {
	// Create services struct for route registration
	services := api.Services{
		TokenService:      tokenService,
		UserService:       userService,
		RoleService:       roleService,
		PositionService:   positionService,
		OperateLogService: logService,
		SettingService:    settingService,
		UserRepository:    userRepository,
	}
	// Pass the gin.Engine to NewWebServer to avoid creating it twice
	return serverpkg.NewWebServer(cfg, engine, services)
}

// ProvideCronManager provides the cron manager.
func ProvideCronManager() *serverpkg.CronManager {
	return serverpkg.NewCronManager()
}

// ProvideHookServer provides the hook server.
func ProvideHookServer() *serverpkg.HookServer {
	return serverpkg.NewHookServer()
}

// ProvideServiceManager provides the service manager with all services.
func ProvideServiceManager(
	cronManager *serverpkg.CronManager,
	webServer *serverpkg.WebServer,
	hookServer *serverpkg.HookServer,
) *task.ServiceManager {
	services := task.NewServiceManager()
	services.AddService(cronManager, webServer, hookServer)
	return services
}

// ============================================================================
// Provider Sets (grouped by functionality)
// ============================================================================

// InfrastructureSet provides all infrastructure dependencies.
var InfrastructureSet = wire.NewSet(
	ProvideConfig,
	ProvideDB,
	ProvideRedis,
)

// RepositorySet provides all repository dependencies.
var RepositorySet = wire.NewSet(
	ProvideUserRepository,
	ProvideRoleRepository,
	ProvideRolePermissionRepository,
	ProvideOperateLogRepository,
	ProvidePositionRepository,
	ProvideServerSettingRepository,
)

// ServiceSet provides all service dependencies.
var ServiceSet = wire.NewSet(
	ProvideTokenService,
	ProvideJwtTokenService,
	ProvideCaptchaService,
	ProvideServerSettingService,
	ProvideOperateLogService,
	ProvidePositionService,
	ProvideRoleService,
	ProvideUserService,
)

// ServerSet provides all HTTP server dependencies.
var ServerSet = wire.NewSet(
	ProvideGinEngine,
	ProvideWebServer,
	ProvideCronManager,
	ProvideHookServer,
	ProvideServiceManager,
)

// AppSet provides all application dependencies.
var AppSet = wire.NewSet(
	InfrastructureSet,
	RepositorySet,
	ServiceSet,
	ServerSet,
)