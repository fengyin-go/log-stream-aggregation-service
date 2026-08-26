package store

import "sync"

// IngestRateSnapshotStore 维护日志摄取速率计数，支持在并发写入期间取快照。
type IngestRateSnapshotStore struct {
	mu    sync.Mutex
	value int
}

func NewIngestRateSnapshotStore(value int) *IngestRateSnapshotStore {
	return &IngestRateSnapshotStore{value: value}
}

// Snapshot 返回采集开始那一刻的计数副本。
//
// 快照必须是不可变的：上层读取快照期间，新日志可能正在写入并调用
// Increment。这里取的是一个冻结的值副本，而不是指向内部字段的指针，
// 因此并发写入不会改变已取得的快照。
func (s *IngestRateSnapshotStore) Snapshot() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.value
}

// Increment 表示一条新日志被写入，摄取速率计数加一。
func (s *IngestRateSnapshotStore) Increment() {
	s.mu.Lock()
	s.value++
	s.mu.Unlock()
}
