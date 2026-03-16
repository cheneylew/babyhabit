package main

import (
	"fmt"
	"os"
	"strings"

	"babyhabit/utils"
)

func main() {
	// 加载环境变量
	envFile, err := os.ReadFile(".env")
	if err != nil {
		fmt.Println("Error reading .env file:", err)
		return
	}

	// 解析环境变量
	lines := strings.Split(string(envFile), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key != "" && !strings.HasPrefix(key, "#") {
				os.Setenv(key, value)
			}
		}
	}

	// 打印环境变量，用于调试
	fmt.Println("Environment variables:")
	fmt.Println("VOICE_APPID:", os.Getenv("VOICE_APPID"))
	fmt.Println("VOICE_ACCESS_TOKEN:", os.Getenv("VOICE_ACCESS_TOKEN"))
	fmt.Println("VOICE_MODEL_ID:", os.Getenv("VOICE_MODEL_ID"))
	fmt.Println("TEXT_API_KEY:", os.Getenv("TEXT_API_KEY"))
	fmt.Println("TEXT_MODEL_ID:", os.Getenv("TEXT_MODEL_ID"))

	// 测试语音合成模型
	fmt.Println("\nTesting speech synthesis model...")
	audioURL, err := utils.GenerateSpeech("hello")
	if err != nil {
		fmt.Println("Error generating speech:", err)
	} else {
		fmt.Println("Speech generated successfully:", audioURL)
	}

	// 测试文本生成模型
	fmt.Println("\nTesting text generation model...")
	example, err := utils.GenerateExampleSentence("hello")
	if err != nil {
		fmt.Println("Error generating example:", err)
	} else {
		fmt.Println("Example sentence generated successfully:", example)
	}
}
