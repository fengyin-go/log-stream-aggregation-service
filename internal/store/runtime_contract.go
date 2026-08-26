package store

import "sync"

type SourceCounterSnapshotStore struct { mu sync.Mutex; value int }
func NewSourceCounterSnapshotStore(value int) *SourceCounterSnapshotStore { return &SourceCounterSnapshotStore{value: value} }
func (s *SourceCounterSnapshotStore) Snapshot() *int { s.mu.Lock(); defer s.mu.Unlock(); return &s.value }
func (s *SourceCounterSnapshotStore) Increment() { s.mu.Lock(); s.value++; s.mu.Unlock() }
