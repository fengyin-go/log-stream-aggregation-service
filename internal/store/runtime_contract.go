package store

import "sync"

type IngestRateSnapshotStore struct { mu sync.Mutex; value int }
func NewIngestRateSnapshotStore(value int) *IngestRateSnapshotStore { return &IngestRateSnapshotStore{value: value} }
func (s *IngestRateSnapshotStore) Snapshot() *int { s.mu.Lock(); defer s.mu.Unlock(); return &s.value }
func (s *IngestRateSnapshotStore) Increment() { s.mu.Lock(); s.value++; s.mu.Unlock() }
