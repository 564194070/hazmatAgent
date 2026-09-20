package llm

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestOpenAILLMClient_Call(t *testing.T) {

	err := godotenv.Load()
	if err != nil {
		t.Fatalf("未找到.env文件，使用系统环境变量，错误%s", err)
	}

	t.Log("启动llmClient链接")

	client := NewOpenAILLMClient()
	resp, err := client.Call("Hello, world!")
	if err != nil {
		t.Fatalf("Call failed: %v", err)
	}
	t.Logf("启动llmClient链接成功,返回结果: %s", resp)
}
