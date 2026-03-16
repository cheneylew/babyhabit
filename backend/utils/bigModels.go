package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// min returns the smaller of x or y
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// getMapKeys returns all keys of a map
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// VoiceSynthesisRequest 语音合成请求结构
type VoiceSynthesisRequest struct {
	Model      string `json:"model"`
	Input      string `json:"input"`
	Parameters struct {
		VoiceType  string `json:"voice_type"`
		Format     string `json:"format"`
		SampleRate int    `json:"sample_rate"`
	} `json:"parameters"`
}

// VoiceSynthesisResponse 语音合成响应结构
type VoiceSynthesisResponse struct {
	Audio string `json:"audio"`
}

// TextGenerationRequest 文本生成请求结构
type TextGenerationRequest struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	Parameters struct {
		MaxTokens   int     `json:"max_tokens"`
		Temperature float64 `json:"temperature"`
	} `json:"parameters"`
}

// TextGenerationResponse 文本生成响应结构
type TextGenerationResponse struct {
	Text string `json:"text"`
}

// generateSignature 生成HMAC-SHA256签名
func generateSignature(secretKey, data string, timestamp int64) string {
	// 构建签名字符串
	signString := fmt.Sprintf("%s%d", data, timestamp)

	// 使用HMAC-SHA256算法生成签名
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(signString))

	// 将签名转换为base64编码
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// GenerateSpeech 生成语音文件
func GenerateSpeech(word string) (string, error) {
	// 从环境变量获取配置
	appID := os.Getenv("VOICE_APPID")
	accessToken := os.Getenv("VOICE_ACCESS_TOKEN")
	secretKey := os.Getenv("VOICE_SECRET_KEY")
	if appID == "" || accessToken == "" || secretKey == "" {
		return "", fmt.Errorf("VOICE_APPID, VOICE_ACCESS_TOKEN, or VOICE_SECRET_KEY not set")
	}

	// 准备请求数据，使用V1 API格式
	reqData := map[string]interface{}{
		"app": map[string]string{
			"appid":   appID,
			"token":   accessToken,
			"cluster": "volcano_tts",
		},
		"user": map[string]string{
			"uid": "123456",
		},
		"audio": map[string]interface{}{
			"voice_type":  "zh_female_wanqudashu_moon_bigtts", // 使用用户指定的音色
			"encoding":    "mp3",
			"speed_ratio": 1.0,
		},
		"request": map[string]interface{}{
			"reqid":     fmt.Sprintf("%d", time.Now().UnixNano()),
			"text":      word,
			"operation": "query",
		},
	}

	data, err := json.Marshal(reqData)
	if err != nil {
		return "", err
	}

	// 打印请求数据，用于调试
	fmt.Println("Request data:", string(data))

	// 发送请求 (V1 API)
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("POST", "https://openspeech.bytedance.com/api/v1/tts", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	// 生成签名
	timestamp := time.Now().Unix()
	signature := generateSignature(secretKey, string(data), timestamp)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer;"+accessToken)
	req.Header.Set("X-Timestamp", fmt.Sprintf("%d", timestamp))
	req.Header.Set("X-Signature", signature)

	resp, err := client.Do(req)
	if err != nil {
		// 处理网络错误，返回空字符串
		fmt.Printf("Speech synthesis network error: %v\n", err)
		return "", nil
	}
	defer resp.Body.Close()

	// 读取整个响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return "", nil
	}

	// 打印响应体，用于调试
	fmt.Println("Response body length:", len(body))
	fmt.Println("First 100 bytes:", string(body[:min(100, len(body))]))

	// 定义audioData变量
	var audioData []byte

	// 尝试解析为JSON
	var respData map[string]interface{}
	if err := json.Unmarshal(body, &respData); err != nil {
		fmt.Printf("Error parsing response as JSON: %v\n", err)

		// 尝试直接作为base64解码
		audioData, err = base64.StdEncoding.DecodeString(string(body))
		if err != nil {
			fmt.Printf("Error decoding response as base64: %v\n", err)
			return "", nil
		}

		// 检查解码后的音频数据是否为空
		if len(audioData) == 0 {
			fmt.Println("Decoded audio data is empty")
			return "", nil
		}
	} else {
		// 从JSON中提取data字段
		if data, ok := respData["data"].(string); ok {
			// 解码音频数据
			audioData, err = base64.StdEncoding.DecodeString(data)
			if err != nil {
				fmt.Printf("Error decoding audio data: %v\n", err)
				return "", nil
			}

			// 检查解码后的音频数据是否为空
			if len(audioData) == 0 {
				fmt.Println("Decoded audio data is empty")
				return "", nil
			}
		} else {
			fmt.Println("Response does not contain data field")
			fmt.Printf("Response keys: %v\n", getMapKeys(respData))
			return "", nil
		}
	}

	// 检查解码后的音频数据是否为空
	if len(audioData) == 0 {
		fmt.Println("Decoded audio data is empty")
		return "", nil
	}

	// 确保目录存在
	wordsDir := filepath.Join("files", "words")
	if err := os.MkdirAll(wordsDir, 0755); err != nil {
		return "", err
	}

	// 生成唯一文件名
	fileName := fmt.Sprintf("%s_%d.mp3", word, time.Now().Unix())
	filePath := filepath.Join(wordsDir, fileName)

	// 写入文件
	if err := os.WriteFile(filePath, audioData, 0644); err != nil {
		return "", err
	}

	// 返回相对路径作为URL
	return "/files/words/" + fileName, nil
}

// GenerateExampleSentence 生成例句
func GenerateExampleSentence(word string) (string, error) {
	// 从环境变量获取配置
	apiKey := os.Getenv("TEXT_API_KEY")
	modelID := os.Getenv("TEXT_MODEL_ID")
	if apiKey == "" || modelID == "" {
		return "", fmt.Errorf("TEXT_API_KEY or TEXT_MODEL_ID not set")
	}

	// 准备请求数据
	prompt := fmt.Sprintf("请为单词 '%s' 生成一个简单的英文例句，不要包含任何解释，只返回例句本身。", word)
	reqData := TextGenerationRequest{
		Model: modelID,
	}
	reqData.Messages = append(reqData.Messages, struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}{Role: "user", Content: prompt})
	reqData.Parameters.MaxTokens = 100
	reqData.Parameters.Temperature = 0.7

	data, err := json.Marshal(reqData)
	if err != nil {
		return "", err
	}

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("POST", "https://ark.cn-beijing.volces.com/api/v3/chat/completions", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed: %s", string(body))
	}

	// 解析响应
	var respData struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", err
	}

	if len(respData.Choices) > 0 {
		return respData.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response from model")
}
