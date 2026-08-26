package store

import "sync"

type ErrorRateSnapshotStore struct { mu sync.Mutex; value int }
func NewErrorRateSnapshotStore(value int) *ErrorRateSnapshotStore { return &ErrorRateSnapshotStore{value: value} }
func (s *ErrorRateSnapshotStore) Snapshot() *int { s.mu.Lock(); defer s.mu.Unlock(); return &s.value }
func (s *ErrorRateSnapshotStore) Increment() { s.mu.Lock(); s.value++; s.mu.Unlock() }
