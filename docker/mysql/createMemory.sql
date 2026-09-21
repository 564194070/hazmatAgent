create database if not exists agent_memory;
use agent_memory;


CREATE TABLE IF NOT EXISTS memory_entry (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '记忆条目唯一ID',
    user_id VARCHAR(36) NOT NULL COMMENT '用户ID',
    session_id VARCHAR(36) NOT NULL COMMENT '会话ID',
    memory_type VARCHAR(32) NOT NULL COMMENT '记忆类型:working / perceptual / episodic / semantic',
    storage_tier VARCHAR(32) NOT NULL DEFAULT 'inContext' COMMENT '存储分层:inContext / longTerm / external',
    content TEXT NOT NULL COMMENT '记忆文本内容',
    create_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间，不会自动更新',
    expire_at TIMESTAMP NULL DEFAULT NULL COMMENT '过期时间，为空代表永久有效',
    metadata JSON COMMENT '扩展元数据:ImportanceScore等',
    INDEX idx_memory_type (memory_type),
    INDEX idx_storage_tier (storage_tier),
    INDEX idx_create_at (create_at),
    INDEX idx_expire_at (expire_at),
    INDEX idx_user_id (user_id),
    INDEX idx_session_id (session_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Agent记忆条目表';
