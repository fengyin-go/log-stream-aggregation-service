package store

import "sync"

type SourceCounterSnapshotStore struct { mu sync.Mutex; value int }
func NewSourceCounterSnapshotStore(value int) *SourceCounterSnapshotStore { return &SourceCounterSnapshotStore{value: value} }
// Snapshot 返回采集开始时固定的来源计数快照。
// 快照必须是冻结的副本：返回指向局部变量的指针，而不是指向 s.value 的指针，
// 否则采集期间 Increment() 的写入会通过指针泄露给上层，使"开始时的计数"变成"最新计数"。
func (s *SourceCounterSnapshotStore) Snapshot() *int {
	s.mu.Lock()
	defer s.mu.Unlock()
	v := s.value
	return &v
}
func (s *SourceCounterSnapshotStore) Increment() { s.mu.Lock(); s.value++; s.mu.Unlock() }
