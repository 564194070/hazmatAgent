package agent

import (
	agentdesc "agentFrame/cmd/xatu-agent/agentDesc"
	"agentFrame/cmd/xatu-agent/llm"
	"agentFrame/cmd/xatu-agent/tools"
	"context"
	"log/slog"
)

type ReActAgent struct {
	llmClient     llm.LLMClientIF
	agentDesc     string
	toolsRegistry tools.ToolRegistryIF
	currentStep   int
	maxSteps      int
	history       []string
}

func NewReActAgent(llmClient llm.LLMClientIF) *ReActAgent {
	return &ReActAgent{llmClient: llmClient, maxSteps: 10, history: make([]string, 0)}
}

func (a *ReActAgent) Run(ctx context.Context, prompt string) (string, error) {

	if a.toolsRegistry == nil {
		a.toolsRegistry = tools.NewToolRegistry()
		a.toolsRegistry.RegisterTool(tools.NewEchoTool())
	}

	// 在限制步数内解决问题
	for a.currentStep = 0; a.currentStep < a.maxSteps; a.currentStep++ {
		reActPromptData := agentdesc.ReActPromptData{
			Tools:    a.toolsRegistry.GetAllToolsDescForPrompt(),
			Question: prompt,
			History:  a.history,
		}

		agentDesc, err := agentdesc.RenderAgentPrompt(reActPromptData)
		if err != nil {
			return "", err
		}

		slog.Info("agentDesc", "agentDesc", agentDesc)
		a.llmClient.Call(agentDesc)
	}

	return "", nil
}
