package store

import "sync"

type AlertContextArchiveStore struct { mu sync.RWMutex; payload []byte }
func NewAlertContextArchiveStore() *AlertContextArchiveStore { return &AlertContextArchiveStore{} }
// Put 将告警上下文写入存档。归档后数据必须与调用方的输入缓冲区
// 彻底分开，因此这里对入参做防御性拷贝，避免调用方复用原切片时
// 把已存档的上下文一并改掉。
func (s *AlertContextArchiveStore) Put(payload []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.payload = append([]byte(nil), payload...)
}

// Snapshot 返回已存档上下文的快照。返回独立拷贝，避免读取方修改
// 返回切片时污染缓存内的存档数据。
func (s *AlertContextArchiveStore) Snapshot() []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]byte(nil), s.payload...)
}
