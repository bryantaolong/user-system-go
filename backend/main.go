package main

import (
	"fmt"
	"os"

	"github.com/bryan/user-system/auth"
	"github.com/bryan/user-system/config"
	"github.com/bryan/user-system/middleware"
	"github.com/bryan/user-system/cache"
	"github.com/bryan/user-system/system"
	"github.com/bryan/user-system/user"
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

	// 5. 初始化模块
	userRepo := auth.NewUserRepository(db)
	roleRepo := auth.NewUserRoleRepository(db)
	profileRepo := user.NewProfileRepository(db)

	// auth 模块
	authSvc := auth.NewService(userRepo, roleRepo, redisSvc, logger)

	// user 模块
	roleSvc := user.NewRoleService(roleRepo)
	userSvc := user.NewService(userRepo, roleSvc, logger)
	fileSvc := user.NewFileService(logger)
	profileSvc := user.NewProfileService(profileRepo, fileSvc, logger)

	// system 模块
	logSvc := system.NewLogService(logger)

	// 6. 初始化 Handler
	authHandler := auth.NewHandler(authSvc, profileSvc.CreateUserProfile)
	userHandler := user.NewHandler(userSvc, profileSvc, authSvc)
	systemHandler := system.NewHandler(logSvc)

	// 7. 初始化 Gin 路由
	r := gin.Default()

	// 全局中间件
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.ErrorHandler(logger))
	r.Use(auth.AuthMiddleware(userRepo, redisSvc))

	// 静态文件服务
	uploadDir := cfg.File.UploadDir
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}
	r.Static("/uploads", uploadDir)

	// 注册路由
	registerRoutes(r, authHandler, userHandler, systemHandler)

	// 8. 启动服务
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
	sqlDB.SetConnMaxLifetime(0)

	logger.Info("数据库连接成功", zap.String("driver", cfg.Driver))
	return db, nil
}

func registerRoutes(
	r *gin.Engine,
	authHandler *auth.Handler,
	userHandler *user.Handler,
	systemHandler *system.Handler,
) {
	// Auth 路由
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.GET("/validate", authHandler.ValidateToken)
		authGroup.GET("/me", auth.RequireAuth(), authHandler.GetCurrentUser)
		authGroup.PUT("/password", auth.RequireAuth(), authHandler.ChangePassword)
		authGroup.DELETE("", auth.RequireAuth(), authHandler.DeleteAccount)
		authGroup.GET("/logout", auth.RequireAuth(), authHandler.Logout)
	}

	// Users 路由
	users := r.Group("/api/users")
	{
		users.POST("", auth.RequireRole("ADMIN"), userHandler.CreateUser)
		users.GET("", auth.RequireRole("ADMIN"), userHandler.ListUsers)
		users.GET("/:userId", auth.RequireRole("ADMIN"), userHandler.GetUserByID)
		users.GET("/username/:username", auth.RequireRole("ADMIN"), userHandler.GetUserByUsername)
		users.POST("/search", auth.RequireRole("ADMIN"), userHandler.QueryUsers)
		users.PUT("/:userId", userHandler.UpdateUser)
		users.PUT("/roles/:userId", auth.RequireRole("ADMIN"), userHandler.ChangeRole)
		users.PUT("/password/:userId", auth.RequireRole("ADMIN"), userHandler.ResetPassword)
		users.PUT("/block/:userId", auth.RequireRole("ADMIN"), userHandler.BlockUser)
		users.PUT("/unblock/:userId", auth.RequireRole("ADMIN"), userHandler.UnblockUser)
		users.DELETE("/:userId", auth.RequireRole("ADMIN"), userHandler.DeleteUser)
	}

	// User Profiles 路由
	profiles := r.Group("/api/user-profiles")
	{
		profiles.POST("/avatar", auth.RequireAuth(), userHandler.UploadAvatar)
		profiles.GET("/:userId", userHandler.GetUserProfileByUserId)
		profiles.GET("/name/:realName", auth.RequireAuth(), userHandler.GetUserProfileByRealName)
		profiles.GET("/me", auth.RequireAuth(), userHandler.GetCurrentUserProfile)
		profiles.PUT("", auth.RequireAuth(), userHandler.UpdateUserProfile)
	}

	// User Roles 路由
	roles := r.Group("/api/user-roles")
	{
		roles.GET("", auth.RequireRole("ADMIN"), userHandler.ListRoles)
	}

	// User Export 路由
	export := r.Group("/api/users/export")
	{
		export.GET("", auth.RequireRole("ADMIN"), userHandler.ExportAllUsers)
	}

	// Admin Logs 路由
	logs := r.Group("/api/admin/logs")
	{
		logs.GET("", auth.RequireRole("ADMIN"), systemHandler.ListLatestLogs)
		logs.GET("/files", auth.RequireRole("ADMIN"), systemHandler.ListLogFiles)
	}
}
