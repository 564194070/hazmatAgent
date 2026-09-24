package controller

import (
	"agentFrame/cmd/xatu-agent/agent"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type ChatController struct {
	agent agent.AgentIF
	mu    sync.Mutex
}

func NewChatController(agent agent.AgentIF) *ChatController {
	return &ChatController{agent: agent}
}

type ChatReq struct {
	UserID    string `json:"userId"`
	SessionID string `json:"sessionId"`
	Message   string `json:"message"`
}

func (h *ChatController) Chat(c *gin.Context) {
	var req ChatReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}

	message := strings.TrimSpace(req.Message)
	//userID := strings.TrimSpace(req.UserID)
	sessionID := strings.TrimSpace(req.SessionID)
	if message == "" || sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "参数错误"})
		return
	}

	uid, ok := c.Get("user_id")
	userID, isUint := uid.(uint)
	if !ok || !isUint || userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "未登录"})
		return
	}

	// ReActAgent 里有步数状态，同一实例不能并发跑
	h.mu.Lock()
	reply, err := h.agent.Run(c.Request.Context(), strconv.FormatUint(uint64(userID), 10), sessionID, message)
	h.mu.Unlock()
	if err != nil {
		slog.Error("智能体执行失败", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "智能体执行失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"reply": reply})
}
