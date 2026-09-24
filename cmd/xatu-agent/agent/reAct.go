package agent

import (
	agentdesc "agentFrame/cmd/xatu-agent/agentDesc"
	"agentFrame/cmd/xatu-agent/llm"
	"agentFrame/cmd/xatu-agent/memory"
	"agentFrame/cmd/xatu-agent/memory/manager"
	"agentFrame/cmd/xatu-agent/tools"
	"context"
	"fmt"
	"log/slog"
	"strings"
)

type ReActAgent struct {
	llmClient     llm.LLMClientIF
	agentDesc     string
	memoryManager manager.MemoryManagerIF
	toolsRegistry tools.ToolRegistryIF
	currentStep   int
	maxSteps      int
}

func NewReActAgent(llmClient llm.LLMClientIF, memoryManager manager.MemoryManagerIF) *ReActAgent {
	return &ReActAgent{llmClient: llmClient, maxSteps: 10, memoryManager: memoryManager}
}

func (a *ReActAgent) SetTools(registry tools.ToolRegistryIF) {
	a.toolsRegistry = registry
}

func (a *ReActAgent) Run(ctx context.Context, userID, sessionID, prompt string) (string, error) {
	if a.toolsRegistry == nil {
		return "", fmt.Errorf("未注册工具")
	}
	if userID == "" || sessionID == "" {
		return "", fmt.Errorf("缺少用户或会话")
	}

	// 问题只写入 Manager 一次，后续 history 全部从 Manager 召回
	if err := a.saveMemory(userID, sessionID, "Question: "+prompt); err != nil {
		return "", err
	}

	reActPromptData := agentdesc.ReActPromptData{
		Tools:    a.toolsRegistry.GetAllToolsDescForPrompt(),
		Question: prompt,
	}

	var lastResp string
	// 在限制步数内解决问题
	for a.currentStep = 0; a.currentStep < a.maxSteps; a.currentStep++ {
		history, err := a.loadHistory(userID, sessionID)
		if err != nil {
			return "", err
		}
		reActPromptData.History = history

		agentDesc, err := agentdesc.RenderAgentPrompt(reActPromptData)
		if err != nil {
			return "", err
		}

		slog.Info("agentDesc", "agentDesc", agentDesc)
		llmResp, err := a.llmClient.Call(agentDesc)
		if err != nil {
			return "", err
		}
		slog.Info("llmResp", "llmResp", llmResp)

		if err := a.saveMemory(userID, sessionID, llmResp); err != nil {
			return "", err
		}

		action, err := parseReActAction(llmResp)
		if err != nil {
			return "", err
		}
		slog.Info("解析Action", "tool", action.Name, "input", action.Input)

		if action.IsFinish() {
			return action.Input, nil
		}

		obs, err := a.executeAction(action)
		if err != nil {
			return "", err
		}
		if err := a.saveMemory(userID, sessionID, "Observation: "+obs); err != nil {
			return "", err
		}
		slog.Info("工具执行结果", "observation", obs)
		lastResp = obs
	}

	return lastResp, nil
}

func (a *ReActAgent) executeAction(action *reActAction) (string, error) {
	tool := a.lookupTool(action.Name)
	if tool == nil {
		return fmt.Sprintf("未知工具 %s，可用工具: %v", action.Name, a.toolsRegistry.ListTools()), nil
	}
	return tool.Execute(parseToolInput(action.Input))
}

func (a *ReActAgent) lookupTool(name string) tools.ToolIF {
	if tool := a.toolsRegistry.GetTool(name); tool != nil {
		return tool
	}
	if tool := a.toolsRegistry.GetTool(name + "工具"); tool != nil {
		return tool
	}
	for _, registered := range a.toolsRegistry.ListTools() {
		if strings.TrimSuffix(registered, "工具") == name {
			return a.toolsRegistry.GetTool(registered)
		}
	}
	return nil
}

func (a *ReActAgent) loadHistory(userID, sessionID string) ([]string, error) {
	entries, err := a.memoryManager.Retrieve(manager.RetrieveOption{
		UserID:    userID,
		SessionID: sessionID,
	})
	if err != nil {
		return nil, err
	}

	history := make([]string, 0, len(entries))
	for _, entry := range entries {
		history = append(history, entry.Content)
	}
	return history, nil
}

func (a *ReActAgent) saveMemory(userID, sessionID, content string) error {
	return a.memoryManager.AddMemory(memory.MemoryEntry{
		UserID:    userID,
		SessionID: sessionID,
		Content:   content,
	})
}
