package memory

import "time"

type MemoryStoreType string

const (
	MemoryStoreTypeInMemory MemoryStoreType = "inMemory" // 内存存储
	MemoryStoreTypeSQLite   MemoryStoreType = "sqlite"   // SQLite存储
	MemoryStoreTypeVectorDB MemoryStoreType = "vectorDB" // 向量存储
	MemoryStoreTypeNeo4j    MemoryStoreType = "neo4j"    // Neo4j存储

	MemoryStoreTypeFile    MemoryStoreType = "file"    // 文件存储
	MemoryStoreTypeRedis   MemoryStoreType = "redis"   // Redis存储
	MemoryStoreTypeMongoDB MemoryStoreType = "mongodb" // MongoDB存储
	MemoryStoreTypeMySQL   MemoryStoreType = "mysql"   // MySQL存储
	MemoryStoreTypeRAG     MemoryStoreType = "rag"     // RAG存储

)

// MemoryType 认知维度：记忆内容类型
type MemoryType string

const (
	MemoryTypeWorking    MemoryType = "working"    // 工作记忆：当前任务中间状态、思考草稿，当前对话上下文
	MemoryTypeEpisodic   MemoryType = "episodic"   // 情景记忆：历史事件、对话、过往任务轨迹，过去交互会话
	MemoryTypeSemantic   MemoryType = "semantic"   // 语义记忆：事实、知识、用户偏好，客观知识，文档，数据
	MemoryTypePerceptual MemoryType = "perceptual" // 感知记忆：原始输入、工具返回原始数据，流程，模板，技能
)

// StorageTier 工程维度：存储分层（对应工程派四类划分）
type StorageTier string

const (
	StorageTierInContext StorageTier = "inContext" // 在上下文窗口内，会话结束丢弃（感知/工作常用）
	StorageTierLongTerm  StorageTier = "longTerm"  // 长期向量库，持久化（episodic/semantic）
	StorageTierExternal  StorageTier = "external"  // 外部独立知识库/文件，大容量静态资料（semantic专用居多）
)

// MemoryEntry 记忆条目结构体，混合双维度。
// 同时作为 GORM 模型：tag 把字段映射到 memory_entry 表，由 ORM 负责 CRUD。
type MemoryEntry struct {
	ID          int            `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	UserID      string         `json:"user_id" gorm:"column:user_id;size:36;not null;index:idx_user_id"`
	SessionID   string         `json:"session_id" gorm:"column:session_id;size:36;not null;index:idx_session_id"`
	MemoryType  MemoryType     `json:"memory_type" gorm:"column:memory_type;size:32;not null;index:idx_memory_type"`
	StorageTier StorageTier    `json:"storage_tier" gorm:"column:storage_tier;size:32;not null;default:inContext;index:idx_storage_tier"`
	Content     string         `json:"content" gorm:"column:content;type:text;not null"`
	CreateAt    time.Time      `json:"create_at" gorm:"column:create_at;autoCreateTime;index:idx_create_at"`
	ExpireAt    *time.Time     `json:"expire_at,omitempty" gorm:"column:expire_at;index:idx_expire_at"`
	Metadata    map[string]any `json:"metadata,omitempty" gorm:"column:metadata;serializer:json;type:json"`
}

func (MemoryEntry) TableName() string {
	return "memory_entry" // 指定数据库表名 memory_entry
}
