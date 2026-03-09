package api

import (
	"babyhabit/controllers"
	"babyhabit/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(router *gin.Engine) {
	// 配置CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 公共路由
	public := router.Group("/api")
	{
		public.POST("/register", controllers.RegisterRequest)
		public.POST("/login", controllers.LoginRequest)
	}

	// 需要认证的路由
	auth := router.Group("/api")
	auth.Use(middleware.JWTAuth())
	{
		// 用户相关
		auth.GET("/user/info", controllers.GetUserInfo)
		auth.PUT("/user/info", controllers.UpdateUserInfo)
		auth.GET("/user/children", controllers.GetChildren)

		// 习惯相关
		auth.GET("/habits", controllers.GetChildHabits)

		// 打卡相关
		auth.POST("/checkin", controllers.Checkin)
		auth.GET("/checkin/records", controllers.GetCheckinRecords)
		auth.GET("/checkin/streak/:habit_id", controllers.GetStreakRecord)

		// 奖励相关
		auth.GET("/rewards", controllers.GetRewardItems)
		auth.POST("/exchange", controllers.Exchange)
		auth.GET("/exchange/records", controllers.GetExchangeRecords)
	}

	// 管理员路由
	admin := router.Group("/api/admin")
	admin.Use(middleware.JWTAuth(), middleware.AdminAuth())
	{
		// 小孩账号管理
		admin.POST("/children", controllers.CreateChild)
		admin.PUT("/children/:id", controllers.UpdateChild)
		admin.DELETE("/children/:id", controllers.DeleteChild)

		// 习惯管理
		admin.POST("/habits", controllers.CreateHabit)
		admin.GET("/habits", controllers.GetHabits)
		admin.GET("/habits/:id", controllers.GetHabit)
		admin.PUT("/habits/:id", controllers.UpdateHabit)
		admin.DELETE("/habits/:id", controllers.DeleteHabit)
		admin.POST("/habits/assign", controllers.AssignHabit)

		// 奖励管理
		admin.POST("/rewards", controllers.CreateRewardItem)
		admin.GET("/rewards", controllers.GetRewardItems)
		admin.GET("/rewards/:id", controllers.GetRewardItem)
		admin.PUT("/rewards/:id", controllers.UpdateRewardItem)
		admin.PUT("/exchange/:id/status", controllers.UpdateExchangeStatus)
	}
}
