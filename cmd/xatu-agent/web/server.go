package web

import (
	"agentFrame/cmd/xatu-agent/utils/tool"
	"agentFrame/cmd/xatu-agent/web/controller"
	"agentFrame/cmd/xatu-agent/web/model"
	"log"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func RunHttpServer() error {
	mysqlUser, err := model.NewMySQLUser()
	if err != nil {
		slog.Error("new mysql user failed", "error", err)
		return err
	}

	userCtl := controller.NewUserController(mysqlUser)

	r := gin.Default()

	// 不需要登录的接口
	public := r.Group("/api/user")
	public.POST("/register", userCtl.Register)
	public.POST("/login", userCtl.Login)

	// 需要JWT登录
	auth := r.Group("/api/user")
	auth.Use(tool.JWTAuth())
	auth.GET("/info", userCtl.GetInfo)

	log.Println("server start :8080")
	return r.Run(":8080")
}
