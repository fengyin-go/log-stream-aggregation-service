package store

import "sync"

type IngestPayloadSnapshotStore struct { mu sync.RWMutex; payload []byte }
func NewIngestPayloadSnapshotStore() *IngestPayloadSnapshotStore { return &IngestPayloadSnapshotStore{} }
func (s *IngestPayloadSnapshotStore) Put(payload []byte) { s.mu.Lock(); defer s.mu.Unlock(); s.payload = payload }
func (s *IngestPayloadSnapshotStore) Snapshot() []byte { s.mu.RLock(); defer s.mu.RUnlock(); return s.payload }
