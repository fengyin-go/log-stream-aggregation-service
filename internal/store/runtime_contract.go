package store

import "sync"

type QueryTokenArchiveStore struct { mu sync.RWMutex; payload []byte }
func NewQueryTokenArchiveStore() *QueryTokenArchiveStore { return &QueryTokenArchiveStore{} }
func (s *QueryTokenArchiveStore) Put(payload []byte) { s.mu.Lock(); defer s.mu.Unlock(); s.payload = payload }
func (s *QueryTokenArchiveStore) Snapshot() []byte { s.mu.RLock(); defer s.mu.RUnlock(); return s.payload }
