package main

import (
	"encoding/json"
	"fmt"
	"github.com/cheneylew/babyhabit/backend/config"
	ollama "github.com/cheneylew/babyhabit/backend/utils/ollama"
	"log"
	"os"
)

func main() {
	// 初始化配置
	if err := config.Init(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// 初始化数据库
	if err := config.InitDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer config.CloseDatabase()

	// 创建 Ollama 客户端
	ollamaClient := ollama.NewOllamaClient("http://localhost:11434", "qwen3.5:0.8b")

	// 获取所有没有中文意思的单词
	words, err := getWordsWithoutChinese()
	if err != nil {
		log.Fatalf("Failed to get words without Chinese: %v", err)
	}

	log.Printf("Found %d words without Chinese meaning", len(words))

	// 逐个处理单词
	for _, word := range words {
		log.Printf("Processing word: %s", word)

		// 生成单词信息
		wordInfo, err := ollamaClient.GenerateWordInfo(word)
		os.Exit(0)
		if err != nil {
			log.Printf("Failed to generate word info for %s: %v", word, err)
			continue
		}

		// 更新数据库
		if err := updateWordInfo(word, wordInfo); err != nil {
			log.Printf("Failed to update word info for %s: %v", word, err)
			continue
		}

		log.Printf("Successfully updated word: %s", word)
	}

	log.Println("Processing completed")
}

// getWordsWithoutChinese 获取所有没有中文意思的单词
func getWordsWithoutChinese() ([]string, error) {
	query := `SELECT english FROM ab_vocabulary WHERE chinese = '' OR chinese IS NULL`
	rows, err := config.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query words: %w", err)
	}
	defer rows.Close()

	var words []string
	for rows.Next() {
		var word string
		if err := rows.Scan(&word); err != nil {
			return nil, fmt.Errorf("failed to scan word: %w", err)
		}
		words = append(words, word)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return words, nil
}

// updateWordInfo 更新单词信息
func updateWordInfo(english string, wordInfo *ollama.WordInfo) error {
	// 转换音标为 JSON 字符串
	phoneticJSON, err := json.Marshal(wordInfo.Phonetic)
	if err != nil {
		return fmt.Errorf("failed to marshal phonetic: %w", err)
	}
	phoneticStr := string(phoneticJSON)

	// 转换例句为 JSON 字符串
	examplesJSON, err := json.Marshal(wordInfo.Examples)
	if err != nil {
		return fmt.Errorf("failed to marshal examples: %w", err)
	}
	examplesStr := string(examplesJSON)

	// 执行更新
	query := `UPDATE ab_vocabulary SET chinese = ?, phonetic = ?, example_sentence = ?, category = ? WHERE english = ?`
	_, err = config.DB.Exec(query, wordInfo.Chinese, phoneticStr, examplesStr, wordInfo.Category, english)
	if err != nil {
		return fmt.Errorf("failed to update word: %w", err)
	}

	return nil
}
