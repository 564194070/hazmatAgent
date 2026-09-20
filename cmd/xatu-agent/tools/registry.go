package tools

import "fmt"

type toolRegistry struct {
	tools map[string]ToolIF
}

func NewToolRegistry() *toolRegistry {
	return &toolRegistry{tools: make(map[string]ToolIF)}
}

func (r *toolRegistry) RegisterTool(tool ToolIF) {
	r.tools[tool.GetName()] = tool
}

func (r *toolRegistry) GetTool(name string) ToolIF {
	return r.tools[name]
}

func (r *toolRegistry) DeleteTool(name string) {
	delete(r.tools, name)
}

func (r *toolRegistry) ListTools() []string {
	toolNames := make([]string, 0, len(r.tools))
	for name := range r.tools {
		toolNames = append(toolNames, name)
	}
	return toolNames
}

func (r *toolRegistry) GetAllToolsDesc() []map[string]string {
	descs := make([]map[string]string, 0, len(r.tools))
	for name, tool := range r.tools {
		descs = append(descs, map[string]string{name: tool.GetDesc()})
	}
	return descs
}

func (r *toolRegistry) GetAllToolsDescForPrompt() []string {
	descs := make([]string, 0, len(r.tools))
	for name, tool := range r.tools {
		descs = append(descs, fmt.Sprintf("%s: %s", name, tool.GetDesc()))
	}
	return descs
}

type ToolRegistryManager struct {
	toolRegistries    map[string]ToolRegistryIF
	enabledRegistries []string
}

func NewToolRegistryManager() *ToolRegistryManager {
	return &ToolRegistryManager{toolRegistries: make(map[string]ToolRegistryIF)}
}
