package store

import (
	"agentFrame/cmd/xatu-agent/memory"
)

// 记忆存储体接口，定义记忆存储体的操作方法 MYSQL/REDIS/Neo4j
type MemoryStorerIF interface {
	// 增加记忆
	Add(entry memory.MemoryEntry) error
	// 根据id删除记忆
	DeleteById(id string) error
	// 根据内容删除记忆
	DeleteByContent(content string) error
	// 更新记忆
	Update(entry memory.MemoryEntry) error
	// 通过user_id获取记忆
	GetByUserId(userId string) ([]memory.MemoryEntry, error)
	// 通过session_id获取记忆
	GetBySessionId(sessionId string) ([]memory.MemoryEntry, error)
	// 通过id获取记忆
	GetById(id string) (memory.MemoryEntry, error)
}
