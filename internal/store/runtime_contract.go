package store

import "sync"

type ErrorRateSnapshotStore struct { mu sync.Mutex; value int }
func NewErrorRateSnapshotStore(value int) *ErrorRateSnapshotStore { return &ErrorRateSnapshotStore{value: value} }
// Snapshot 返回计数在调用时刻的不可变副本：先拷贝到局部变量再取地址，
// 使每次快照都指向独立的新分配，后续 Increment 修改 s.value 不会影响已交出的快照。
func (s *ErrorRateSnapshotStore) Snapshot() *int { s.mu.Lock(); defer s.mu.Unlock(); v := s.value; return &v }
func (s *ErrorRateSnapshotStore) Increment() { s.mu.Lock(); s.value++; s.mu.Unlock() }
