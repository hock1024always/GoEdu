package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// DeepSeekRequest 是请求体的结构
type DeepSeekRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// Message 是消息的结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func GetAIResponse(prompt string) (string, error) {
	// DeepSeek API 的 URL 和 API 密钥
	apiURL := "https://openrouter.ai/api/v1/chat/completions"                             // 替换为正确的 API URL
	apiKey := "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" // 替换为你的 DeepSeek API 密钥

	// 构造请求体
	requestBody := DeepSeekRequest{
		Model: "deepseek/deepseek-r1-distill-llama-70b:free", // 使用 deepseek 模型（或其他支持的模型）
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
	}

	// 将请求体序列化为 JSON
	jsonData, err := json.Marshal(requestBody)
	fmt.Println(string(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to %s: %v", apiURL, err)
	}
	defer resp.Body.Close()

	// 检查 HTTP 状态码
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("API request to %s failed with status: %d, response: %s", apiURL, resp.StatusCode, string(body))
	}

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	if err != nil {
		return "", fmt.Errorf("failed to read response body from %s: %v", apiURL, err)
	}

	// 解析响应体
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response from %s: %v", apiURL, err)
	}

	// 提取 AI 的回答
	if choices, ok := response["choices"].([]interface{}); ok && len(choices) > 0 {
		if message, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{}); ok {
			if content, ok := message["content"].(string); ok {
				return content, nil
			}
		}
	}
	return "", fmt.Errorf("failed to extract response content from %s", apiURL)
}
