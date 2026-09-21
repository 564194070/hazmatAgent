package manager

import (
	"agentFrame/cmd/xatu-agent/memory"
	"agentFrame/cmd/xatu-agent/memory/store"
)

// 开关，管理记忆存储体实际存储支撑
type MemoryStoreManager struct {
	memoryStores  map[memory.MemoryStoreType]store.MemoryStorerIF
	config        map[memory.MemoryStoreType]any
	enabledStores []memory.MemoryStoreType
}

func NewMemoryStoreManager(enabledStores []memory.MemoryStoreType, config map[memory.MemoryStoreType]any) (*MemoryStoreManager, error) {
	m := &MemoryStoreManager{
		memoryStores:  make(map[memory.MemoryStoreType]store.MemoryStorerIF),
		enabledStores: enabledStores,
		config:        config,
	}

	for _, storeType := range enabledStores {
		switch storeType {
		case memory.MemoryStoreTypeMySQL:
			store, err := store.NewMySQLStore()
			if err != nil {
				return nil, err
			}
			m.memoryStores[storeType] = store
		}
	}

	return m, nil
}

func (m *MemoryStoreManager) AddMemory(entry memory.MemoryEntry) error {

	for _, store := range m.memoryStores {
		err := store.Add(entry)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MemoryStoreManager) Retrieve(opt RetrieveOption) ([]*memory.MemoryEntry, error) {

	var entries []memory.MemoryEntry
	var err error

	for _, store := range m.memoryStores {
		entries, err = store.GetByUserId(opt.UserID)
		if err != nil {
			return nil, err
		}
	}

	result := make([]*memory.MemoryEntry, 0, len(entries))
	for i := range entries {
		result = append(result, &entries[i])
	}
	return result, nil
}

func (m *MemoryStoreManager) CompressMemories(userID string, rawEntries []*memory.MemoryEntry) (*memory.MemoryEntry, error) {
	// 记忆压缩，待实现
	return nil, nil
}

func (m *MemoryStoreManager) UpdateImportance(entryID string) error {

	// 记忆重要性更新，待实现
	return nil
}

func (m *MemoryStoreManager) CleanExpired(now int64) (int64, error) {

	// 清理过期记忆，待实现
	return 0, nil
}
