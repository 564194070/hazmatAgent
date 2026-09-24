package manager

import (
	"agentFrame/cmd/xatu-agent/memory"
)

type RetrieveOption struct {
	UserID      string
	SessionID   string
	Query       string
	TopK        int
	MaxToken    int // 新增：召回记忆总token上限
	MemoryTypes []memory.MemoryType
}

// 记忆管理器接口
type MemoryManagerIF interface {
	// 内部：调用store.Save，同时可自动更新importance
	AddMemory(entry memory.MemoryEntry) error

	// 召回记忆：混合检索（结构化过滤 + 向量语义召回）
	// 输入用户当前query，返回排序好的相关记忆，给Agent拼入prompt
	Retrieve(opt RetrieveOption) ([]*memory.MemoryEntry, error)

	// 批量压缩：把多条raw perceptual记忆，压缩生成episodic记忆
	// 压缩完成后，可选择标记旧感知记忆为过期
	CompressMemories(userID string, rawEntries []*memory.MemoryEntry) (*memory.MemoryEntry, error)

	// 更新记忆重要性分数（遗忘逻辑）
	UpdateImportance(entryID string) error

	// 清理过期记忆（代理调用底层store.DeleteExpired）
	CleanExpired(now int64) (int64, error)
}
