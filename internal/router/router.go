package router

import (
	"github.com/bryantaolong/system/internal/handler"
	"github.com/bryantaolong/system/internal/middleware"
	"github.com/bryantaolong/system/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(
	db *gorm.DB,
	authService *service.AuthService,
	userService *service.UserService,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	// 公开接口
	public := r.Group("/api/auth")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
		public.GET("/validate", authHandler.Validate)
	}

	// 受保护接口
	protected := r.Group("/api")
	protected.Use(middleware.AuthRequired(authService))
	{
		protected.GET("/auth/me", authHandler.Me)
		protected.GET("/auth/logout", authHandler.Logout)

		admin := protected.Group("/user")
		admin.Use(middleware.RoleRequired(authService, "ROLE_ADMIN"))
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
		}
	}

	return r
}
