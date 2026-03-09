package main

import (
	"babyhabit/api"
	"babyhabit/config"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	if err := config.Init(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// 初始化数据库（暂时注释，以便测试服务器启动）
	if err := config.InitDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer config.CloseDatabase()
	// 创建Gin引擎
	router := gin.Default()

	// 设置路由
	api.SetupRoutes(router)

	// 启动服务器
	port := config.AppConfig.Server.Port
	log.Printf("Server starting on port %s...", port)
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
