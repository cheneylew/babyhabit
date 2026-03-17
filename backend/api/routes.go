package api

import (
	"babyhabit/controllers"
	"babyhabit/middleware"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(router *gin.Engine) {
	// 配置 CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 静态文件服务（前端打包后的文件）
	// 服务 /assets/ 等静态资源
	router.Static("/assets", "./static/assets")

	// 服务上传的文件（照片等）
	router.Static("/files", "./files")

	// 前端页面路由（支持 HTML5 History 模式）
	router.NoRoute(func(c *gin.Context) {
		// API 路由不处理
		if len(c.Request.URL.Path) > 4 && c.Request.URL.Path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}

		// 返回前端首页
		c.File("./static/index.html")
	})

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
		auth.POST("/checkin/submit-rate", controllers.SubmitSelfRate)

		// 打卡回退（管理员）
		auth.POST("/admin/checkin/rollback", controllers.RollbackCheckin)
		auth.GET("/admin/child/checkin-records", controllers.GetChildCheckinRecords)

		// 奖励相关
		auth.GET("/rewards", controllers.GetRewardItems)
		auth.POST("/exchange", controllers.Exchange)
		auth.GET("/exchange/records", controllers.GetExchangeRecords)

		// 名言警句API
		auth.GET("/quote/random", controllers.GetRandomQuote)

		// 词汇学习API
		auth.GET("/vocabulary/plan", controllers.GetVocabularyPlan)
		auth.GET("/vocabulary/start", controllers.StartVocabularyLearning)
		auth.GET("/vocabulary/options", controllers.GetVocabularyOptions)
		auth.POST("/vocabulary/record", controllers.RecordVocabularyLearning)
		auth.GET("/vocabulary/stats", controllers.GetVocabularyStats)
		// 生成例句音频
		auth.POST("/vocabulary/generate-sentence-audio", controllers.GenerateSentenceAudio)
	}

	// 管理员路由
	admin := router.Group("/api/admin")
	admin.Use(middleware.JWTAuth(), middleware.AdminAuth())
	{
		// 小孩账号管理
		admin.GET("/children", controllers.GetChildren)
		admin.GET("/children/:id", controllers.GetChild)
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
		admin.POST("/habits/batch-assign", controllers.BatchAssignHabits)
		admin.GET("/assigned-habits", controllers.GetAssignedHabits)
		admin.DELETE("/habit-assignments/:id", controllers.DeleteHabitAssignment)

		// 奖励管理
		admin.POST("/rewards", controllers.CreateRewardItem)
		admin.GET("/rewards", controllers.GetRewardItems)
		admin.GET("/rewards/:id", controllers.GetRewardItem)
		admin.PUT("/rewards/:id", controllers.UpdateRewardItem)
		admin.DELETE("/rewards/:id", controllers.DeleteRewardItem)
		admin.GET("/exchanges", controllers.GetAllExchangeRecords)
		admin.PUT("/exchange/:id/status", controllers.UpdateExchangeStatus)
		admin.DELETE("/exchanges/:id", controllers.DeleteExchangeRecord)

		// 名言警句管理
		admin.POST("/quotes", controllers.CreateQuote)
		admin.GET("/quotes", controllers.GetQuotes)
		admin.GET("/quotes/:id", controllers.GetQuote)
		admin.PUT("/quotes/:id", controllers.UpdateQuote)
		admin.DELETE("/quotes/:id", controllers.DeleteQuote)
		admin.POST("/quotes/batch", controllers.BatchCreateQuotes)
		admin.DELETE("/quotes/batch", controllers.BatchDeleteQuotes)

		// 词汇学习管理
		admin.POST("/vocabulary", controllers.CreateVocabulary)
		admin.GET("/vocabulary", controllers.GetVocabularies)
		admin.GET("/vocabulary/:id", controllers.GetVocabulary)
		admin.PUT("/vocabulary/:id", controllers.UpdateVocabulary)
		admin.DELETE("/vocabulary/:id", controllers.DeleteVocabulary)
		admin.POST("/vocabulary/batch", controllers.BatchCreateVocabulary)
		admin.DELETE("/vocabulary/batch", controllers.BatchDeleteVocabulary)
	}
}
