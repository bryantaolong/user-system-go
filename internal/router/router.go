package router

import (
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/bryantaolong/system/internal/handler"
	"github.com/bryantaolong/system/internal/middleware"
	"github.com/bryantaolong/system/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	redisClient *redis.Client,
	authService *service.AuthService,
	userService *service.UserService,
	userRoleService *service.UserRoleService,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// CORS 配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	userRoleHandler := handler.NewUserRoleHandler(userRoleService)

	// 公开接口
	public := r.Group("/api/auth")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
		public.GET("/validate", authHandler.Validate)
	}

	// 受保护接口
	protected := r.Group("/api")
	protected.Use(middleware.AuthRequired(redisClient))
	{
		protected.GET("/auth/me", authHandler.Me)
		protected.GET("/auth/logout", authHandler.Logout)

		admin := protected.Group("/user")
		admin.Use(middleware.RoleRequired("ROLE_ADMIN"))
		{
			admin.POST("/all", userHandler.GetAllUsers)
			admin.GET("/:userId", userHandler.GetUserByID)
			admin.GET("/username/:username", userHandler.GetUserByUsername)
			admin.POST("/search", userHandler.SearchUsers)
			admin.PUT("/:userId", userHandler.UpdateUser)
			admin.PUT("/:userId/role", userHandler.ChangeRole)
			admin.PUT("/:userId/password", userHandler.ChangePassword)
			admin.PUT("/:userId/password/force/:newPassword", userHandler.ChangePasswordForcefully)
			admin.PUT("/:userId/block", userHandler.BlockUser)
			admin.PUT("/:userId/unblock", userHandler.UnblockUser)
			admin.DELETE("/:userId", userHandler.DeleteUser)

			admin.GET("/role/all", userRoleHandler.ListRoles)
		}
	}

	return r
}
