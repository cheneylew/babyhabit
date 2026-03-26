package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OllamaClient Ollama API 客户端
type OllamaClient struct {
	BaseURL string
	Model   string
}

// NewOllamaClient 创建新的 Ollama 客户端
func NewOllamaClient(baseURL, model string) *OllamaClient {
	return &OllamaClient{
		BaseURL: baseURL,
		Model:   model,
	}
}

// GenerateWordInfoRequest 生成单词信息请求
type GenerateWordInfoRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Message   struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Done bool `json:"done"`
}

// WordInfo 单词信息结构
type WordInfo struct {
	English  string       `json:"english"`
	Chinese  string       `json:"chinese"`
	Phonetic PhoneticInfo `json:"phonetic"`
	Examples []Example    `json:"examples"`
	Category string       `json:"category"`
	AudioURL string       `json:"audio_url"`
}

// PhoneticInfo 音标信息结构
type PhoneticInfo struct {
	UK string `json:"uk"`
	US string `json:"us"`
}

// Example 例句结构
type Example struct {
	English string `json:"english"`
	Chinese string `json:"chinese"`
}

// Chat 调用 Ollama API 进行聊天
func (c *OllamaClient) Chat(prompt string) (string, error) {
	// 构建请求
	reqBody, err := json.Marshal(map[string]interface{}{
		"model": c.Model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"think":  false,
		"stream": false,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// 发送请求
	client := &http.Client{
		Timeout: 600 * time.Second,
	}
	resp, err := client.Post(c.BaseURL+"/api/chat", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// 解析响应
	var ollamaResp ChatResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return ollamaResp.Message.Content, nil
}

// GenerateWordInfo 使用 Ollama 生成单词信息
func (c *OllamaClient) GenerateWordInfo(word string) (*WordInfo, error) {
	// 构建提示词
	prompt := fmt.Sprintf(`请为单词 "%s" 生成以下信息：
1. 中文意思
2. 英音和美音音标
3. 3个例句（每个例句包含英文和中文翻译）
4. 单词类别（如：名词、动词、形容词等）

请以 JSON 格式返回，结构如下：
{
  "chinese": "中文意思",
  "phonetic": {
    "uk": "英音音标",
    "us": "美音音标"
  },
  "examples": [
    {"english": "例句1", "chinese": "翻译1"},
    {"english": "例句2", "chinese": "翻译2"},
    {"english": "例句3", "chinese": "翻译3"}
  ],
  "category": "单词类别"
}`, word)
	fmt.Println(prompt)

	// 调用 Chat 函数获取响应
	response, err := c.Chat(prompt)
	if err != nil {
		return nil, err
	}

	// 打印响应内容
	fmt.Println(response)

	// 解析返回的 JSON 数据
	var wordInfo WordInfo
	wordInfo.English = word
	if err := json.Unmarshal([]byte(response), &wordInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal word info: %w", err)
	}

	return &wordInfo, nil
}
