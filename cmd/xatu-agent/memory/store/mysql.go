package store

import (
	"agentFrame/cmd/xatu-agent/memory"
	"agentFrame/cmd/xatu-agent/utils/client"
	"database/sql"
	"log/slog"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLStore struct {
	sqlDB  *sql.DB
	gormDB *gorm.DB
}

func NewMySQLStore() (*MySQLStore, error) {
	config := client.NewMySQLClientConfig(os.Getenv("MYSQL_MEMORY_DBNAME"))
	dsn := config.GetDSN()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("open db failed", "error", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("get sql db failed", "error", err)
		return nil, err
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量。
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	if err := sqlDB.Ping(); err != nil {
		slog.Error("ping db failed", "error", err)
		_ = sqlDB.Close()
		return nil, err
	}

	return &MySQLStore{sqlDB: sqlDB, gormDB: db}, nil
}

func (s *MySQLStore) Add(entry memory.MemoryEntry) error {
	res := s.gormDB.Create(&entry)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (s *MySQLStore) DeleteById(id string) error {
	res := s.gormDB.Delete(&memory.MemoryEntry{}, id)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (s *MySQLStore) DeleteByContent(content string) error {
	res := s.gormDB.Where("content = ?", content).Delete(&memory.MemoryEntry{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (s *MySQLStore) Update(entry memory.MemoryEntry) error {
	res := s.gormDB.Save(&entry)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (s *MySQLStore) GetById(id string) (memory.MemoryEntry, error) {
	var entry memory.MemoryEntry
	res := s.gormDB.Where("id = ?", id).First(&entry)
	if res.Error != nil {
		return memory.MemoryEntry{}, res.Error
	}
	return entry, nil
}

func (s *MySQLStore) GetByUserId(userId string) ([]memory.MemoryEntry, error) {
	var entries []memory.MemoryEntry
	res := s.gormDB.Where("user_id = ?", userId).Find(&entries)
	if res.Error != nil {
		return nil, res.Error
	}
	return entries, nil
}

func (s *MySQLStore) GetBySessionId(sessionId string) ([]memory.MemoryEntry, error) {
	var entries []memory.MemoryEntry
	res := s.gormDB.Where("session_id = ?", sessionId).Find(&entries)
	if res.Error != nil {
		return nil, res.Error
	}
	return entries, nil
}

func (s *MySQLStore) GetByUserAndSession(userId, sessionId string) ([]memory.MemoryEntry, error) {
	var entries []memory.MemoryEntry
	res := s.gormDB.Where("user_id = ? AND session_id = ?", userId, sessionId).
		Order("create_at asc").
		Find(&entries)
	if res.Error != nil {
		return nil, res.Error
	}
	return entries, nil
}
