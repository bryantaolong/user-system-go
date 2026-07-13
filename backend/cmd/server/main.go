package main

import (
	"fmt"
	"os"

	"github.com/bryan/user-system/internal/config"
	"github.com/bryan/user-system/internal/handler"
	"github.com/bryan/user-system/internal/middleware"
	"github.com/bryan/user-system/internal/pkg/redis"
	"github.com/bryan/user-system/internal/repository"
	"github.com/bryan/user-system/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. 加载配置
	cfg := config.Load()

	// 2. 初始化日志
	logger := config.InitLogger(&cfg.Logging)
	defer logger.Sync()

	// 3. 初始化数据库
	db, err := initDB(&cfg.Database, logger)
	if err != nil {
		logger.Fatal("数据库连接失败", zap.Error(err))
	}

	// 4. 初始化 Redis
	redisSvc := redis.NewRedisClient(&cfg.Redis, logger)

	// 5. 初始化 Repository
	userRepo := repository.NewUserRepository(db)
	userProfileRepo := repository.NewUserProfileRepository(db)
	userRoleRepo := repository.NewUserRoleRepository(db)

	// 6. 初始化 Service
	userRoleSvc := service.NewUserRoleService(userRoleRepo)
	localFileSvc := service.NewLocalFileService(logger)
	userProfileSvc := service.NewUserProfileService(userProfileRepo, localFileSvc, logger)
	userSvc := service.NewUserService(userRepo, userRoleSvc, logger)
	authSvc := service.NewAuthService(userRepo, userRoleRepo, redisSvc, logger)
	logSvc := service.NewLogService(logger)

	// 7. 初始化 Handler
	authHandler := handler.NewAuthHandler(authSvc, userProfileSvc)
	userHandler := handler.NewUserHandler(userSvc, userProfileSvc)
	userProfileHandler := handler.NewUserProfileHandler(userProfileSvc, userSvc, authSvc)
	userRoleHandler := handler.NewUserRoleHandler(userRoleSvc)
	systemLogHandler := handler.NewSystemLogHandler(logSvc)

	// 8. 初始化 Gin 路由
	r := gin.Default()

	// 全局中间件
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.ErrorHandler(logger))
	r.Use(middleware.AuthMiddleware(userRepo, redisSvc))

	// 静态文件服务
	uploadDir := cfg.File.UploadDir
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}
	r.Static("/uploads", uploadDir)

	// 注册路由
	registerRoutes(r, authHandler, userHandler, userProfileHandler, userRoleHandler, systemLogHandler)

	// 9. 启动服务
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info("服务启动", zap.String("addr", addr))
	if err := r.Run(addr); err != nil {
		logger.Fatal("服务启动失败", zap.Error(err))
	}
}

func initDB(cfg *config.DatabaseConfig, logger *zap.Logger) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.Driver {
	case "postgres":
		dialector = postgres.Open(cfg.DSN())
	case "mysql":
		dialector = mysql.Open(cfg.DSN())
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(0) // 使用配置

	logger.Info("数据库连接成功", zap.String("driver", cfg.Driver))
	return db, nil
}

func registerRoutes(
	r *gin.Engine,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	userProfileHandler *handler.UserProfileHandler,
	userRoleHandler *handler.UserRoleHandler,
	systemLogHandler *handler.SystemLogHandler,
) {
	// Auth 路由
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.GET("/validate", authHandler.ValidateToken)
		auth.GET("/me", middleware.RequireAuth(), authHandler.GetCurrentUser)
		auth.PUT("/password", middleware.RequireAuth(), authHandler.ChangePassword)
		auth.DELETE("", middleware.RequireAuth(), authHandler.DeleteAccount)
		auth.GET("/logout", middleware.RequireAuth(), authHandler.Logout)
	}

	// Users 路由
	users := r.Group("/api/users")
	{
		users.POST("", middleware.RequireRole("ADMIN"), userHandler.CreateUser)
		users.GET("", middleware.RequireRole("ADMIN"), userHandler.ListUsers)
		users.GET("/:userId", middleware.RequireRole("ADMIN"), userHandler.GetUserByID)
		users.GET("/username/:username", middleware.RequireRole("ADMIN"), userHandler.GetUserByUsername)
		users.POST("/search", middleware.RequireRole("ADMIN"), userHandler.QueryUsers)
		users.PUT("/:userId", userHandler.UpdateUser) // 权限在handler中判断
		users.PUT("/roles/:userId", middleware.RequireRole("ADMIN"), userHandler.ChangeRole)
		users.PUT("/password/:userId", middleware.RequireRole("ADMIN"), userHandler.ResetPassword)
		users.PUT("/block/:userId", middleware.RequireRole("ADMIN"), userHandler.BlockUser)
		users.PUT("/unblock/:userId", middleware.RequireRole("ADMIN"), userHandler.UnblockUser)
		users.DELETE("/:userId", middleware.RequireRole("ADMIN"), userHandler.DeleteUser)
	}

	// User Profiles 路由
	profiles := r.Group("/api/user-profiles")
	{
		profiles.POST("/avatar", middleware.RequireAuth(), userProfileHandler.UploadAvatar)
		profiles.GET("/:userId", userProfileHandler.GetUserProfileByUserId)
		profiles.GET("/name/:realName", middleware.RequireAuth(), userProfileHandler.GetUserProfileByRealName)
		profiles.GET("/me", middleware.RequireAuth(), userProfileHandler.GetCurrentUserProfile)
		profiles.PUT("", middleware.RequireAuth(), userProfileHandler.UpdateUserProfile)
	}

	// User Roles 路由
	roles := r.Group("/api/user-roles")
	{
		roles.GET("", middleware.RequireRole("ADMIN"), userRoleHandler.ListRoles)
	}

	// User Export 路由
	export := r.Group("/api/users/export")
	{
		export.GET("", middleware.RequireRole("ADMIN"), userHandler.ExportAllUsers)
	}

	// Admin Logs 路由
	logs := r.Group("/api/admin/logs")
	{
		logs.GET("", middleware.RequireRole("ADMIN"), systemLogHandler.ListLatestLogs)
		logs.GET("/files", middleware.RequireRole("ADMIN"), systemLogHandler.ListLogFiles)
	}
}
