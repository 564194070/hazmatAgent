package agent

import (
	"encoding/json"
	"fmt"
	"strings"
)

type reActAction struct {
	Name  string
	Input string
}

func (a *reActAction) IsFinish() bool {
	return strings.EqualFold(a.Name, "Finish")
}

// parseReActAction 从 llmResp 中取出 Action: ToolName[ToolInput]。
// 用 LastIndex 找闭合 ]，这样 ToolInput 里如果带括号也不会截断。
func parseReActAction(llmResp string) (*reActAction, error) {
	const prefix = "Action:"
	idx := strings.Index(llmResp, prefix)
	if idx < 0 {
		idx = strings.Index(llmResp, "action:")
		if idx < 0 {
			return nil, fmt.Errorf("llmResp 中没有 Action 字段")
		}
	}

	rest := strings.TrimSpace(llmResp[idx+len(prefix):])
	open := strings.Index(rest, "[")
	if open < 0 {
		return nil, fmt.Errorf("Action 不是 ToolName[ToolInput] 格式: %q", rest)
	}
	close := strings.LastIndex(rest, "]")
	if close <= open {
		return nil, fmt.Errorf("Action 缺少闭合括号: %q", rest)
	}

	name := strings.TrimSpace(rest[:open])
	if name == "" {
		return nil, fmt.Errorf("Action 缺少工具名")
	}
	return &reActAction{
		Name:  name,
		Input: rest[open+1 : close],
	}, nil
}

func parseToolInput(raw string) map[string]any {
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, "{") && strings.HasSuffix(raw, "}") {
		var m map[string]any
		if err := json.Unmarshal([]byte(raw), &m); err == nil {
			return m
		}
	}
	return map[string]any{
		"path":  raw,
		"input": raw,
	}
}
