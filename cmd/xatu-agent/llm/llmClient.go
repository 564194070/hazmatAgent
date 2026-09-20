package llm

import (
	"context"
	"log/slog"
	"os"

	openaiSDK "github.com/sashabaranov/go-openai"
)

type LLMClient interface {
	Call(prompt string) (string, error)
}

type OpenAILLMClient struct {
	apiKey    string
	model     string
	baseURL   string
	timeout   int
	maxTokens int
}

func NewOpenAILLMClient() *OpenAILLMClient {
	return &OpenAILLMClient{
		apiKey:    os.Getenv("LLM_API_KEY"),
		model:     os.Getenv("LLM_MODEL_ID"),
		baseURL:   os.Getenv("LLM_BASE_URL"),
		timeout:   10,
		maxTokens: 1000,
	}
}

func (c *OpenAILLMClient) Call(prompt string) (string, error) {

	slog.Info("调用大模型", "apiKey", c.apiKey, "model", c.model, "baseURL", c.baseURL)
	config := openaiSDK.DefaultConfig(c.apiKey)
	config.BaseURL = c.baseURL
	client := openaiSDK.NewClientWithConfig(config)

	// 顶层对话包装
	req := openaiSDK.ChatCompletionRequest{
		Model: c.model,
		// 对话里面的一句话
		Messages: []openaiSDK.ChatCompletionMessage{
			{
				// 用户角色
				/*
					"system"：系统提示词，设定模型角色、规则
					"user"：用户输入
					"assistant"：模型输出回复
					"tool"：工具调用返回结果（工具返回消息）
					"function"：旧版 function 角色（已经被 tool_calls 替代，兼容老协议）
				*/
				Role:    openaiSDK.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	}

	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		slog.Error("大模型调用错误", "错误信息", err)
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}
