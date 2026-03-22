package main

import (
	"babyhabit/api"
	"babyhabit/config"
	"babyhabit/utils"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	if err := config.Init(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// 测试chat接口
	test()

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

func test() {
	testChatStream()
	os.Exit(0)
}

func testChat() {
	prompt := "你好"
	response, err := utils.Chat(prompt)
	if err != nil {
		log.Fatalf("Failed to chat: %v", err)
	}
	log.Printf("Chat response: %s", response)
}

func testChatStream() {
	prompt := "fabricate怎么学好这个单词？"
	err := utils.ChatStream(prompt, func(chunk string) bool {
		log.Printf("Chat stream chunk: %s", chunk)
		return true
	}, nil)
	if err != nil {
		log.Fatalf("Failed to chat stream: %v", err)
	}
}
