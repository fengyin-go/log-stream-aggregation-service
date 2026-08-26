package store

import "sync"

type AlertContextArchiveStore struct { mu sync.RWMutex; payload []byte }
func NewAlertContextArchiveStore() *AlertContextArchiveStore { return &AlertContextArchiveStore{} }
func (s *AlertContextArchiveStore) Put(payload []byte) { s.mu.Lock(); defer s.mu.Unlock(); s.payload = payload }
func (s *AlertContextArchiveStore) Snapshot() []byte { s.mu.RLock(); defer s.mu.RUnlock(); return s.payload }
