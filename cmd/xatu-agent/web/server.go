package web

import (
	"agentFrame/cmd/xatu-agent/agent"
	"agentFrame/cmd/xatu-agent/llm"
	"agentFrame/cmd/xatu-agent/memory"
	"agentFrame/cmd/xatu-agent/memory/manager"
	"agentFrame/cmd/xatu-agent/tools"
	"agentFrame/cmd/xatu-agent/utils/tool"
	"agentFrame/cmd/xatu-agent/web/controller"
	"agentFrame/cmd/xatu-agent/web/model"
	"log"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func RunHttpServer() error {

	// 用户管理相关控制器
	mysqlUser, err := model.NewMySQLUser()
	if err != nil {
		slog.Error("new mysql user failed", "error", err)
		return err
	}

	userCtl := controller.NewUserController(mysqlUser)

	// ReActAgent 智能体相关控制器
	memoryManager, err := manager.NewMemoryStoreManager([]memory.MemoryStoreType{memory.MemoryStoreTypeMySQL}, nil)
	if err != nil {
		slog.Error("new memory manager failed", "error", err)
		return err
	}
	reActAgent := agent.NewReActAgent(llm.NewOpenAILLMClient(), memoryManager)
	registry := tools.NewToolRegistry()
	registry.RegisterTool(tools.NewEchoTool())
	reActAgent.SetTools(registry)
	chatCtl := controller.NewChatController(reActAgent)

	// 登录逻辑
	r := gin.Default()

	// 不需要登录的接口
	public := r.Group("/api/user")
	public.POST("/register", userCtl.Register)
	public.POST("/login", userCtl.Login)

	// 需要JWT登录
	auth := r.Group("/api")
	auth.Use(tool.JWTAuth())
	auth.GET("/user/info", userCtl.GetInfo)
	auth.POST("/chat", chatCtl.Chat)

	log.Println("server start :8080")
	return r.Run(":8080")
}
