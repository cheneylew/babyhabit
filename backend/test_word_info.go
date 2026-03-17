package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"

	"babyhabit/utils"
)

func main() {
	// 加载.env文件
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found or cannot be loaded")
	}

	// 测试词汇
	word := "apple"

	// 生成完整的单词信息
	wordInfo, err := utils.GenerateWordInfo(word)
	if err != nil {
		log.Fatalf("Error generating word info: %v", err)
	}

	// 打印结果
	fmt.Printf("Word: %s\n", word)
	fmt.Printf("Chinese: %s\n", wordInfo.Chinese)
	fmt.Printf("Phonetic: %s\n", wordInfo.Phonetic)
	fmt.Println("Examples:")
	for i, example := range wordInfo.Examples {
		fmt.Printf("  %d. %s\n", i+1, example)
	}
	fmt.Printf("Category: %s\n", wordInfo.Category)
	fmt.Printf("Audio URL: %s\n", wordInfo.AudioURL)

	fmt.Println("Test completed successfully!")
}
