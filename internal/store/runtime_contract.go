package store

import "sync"

type QueryTokenArchiveStore struct {
	mu      sync.RWMutex
	payload []byte
}

func NewQueryTokenArchiveStore() *QueryTokenArchiveStore { return &QueryTokenArchiveStore{} }

// Put 防御性拷贝 payload：归档的是独立副本，调用方之后复用同一底层数组
// 写入新令牌也不会覆盖已归档的内容。
func (s *QueryTokenArchiveStore) Put(payload []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.payload = append([]byte(nil), payload...)
}

// Snapshot 返回独立副本，避免外部修改回写内部缓存。
func (s *QueryTokenArchiveStore) Snapshot() []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]byte(nil), s.payload...)
}
