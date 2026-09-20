package tools

type ToolIF interface {
	Execute(map[string]any) (string, error)
	GetDesc() string
	GetName() string
}

// 不同的工具，可以用自己专属的工具仓库
type ToolRegistryIF interface {
	RegisterTool(tool ToolIF)
	GetTool(name string) ToolIF
	DeleteTool(name string)
	// 展示工具列表
	ListTools() []string
	// 工具名称 -> 工具功能
	GetAllToolsDesc() []map[string]string
	// 展示 prompt 需要的工具列表
	GetAllToolsDescForPrompt() []string
}
