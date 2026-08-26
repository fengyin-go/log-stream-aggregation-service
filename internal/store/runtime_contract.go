package store

import (
    "fmt"
    "sync"
)

type SourceReconnectRetryTemporaryError struct { Key string }
func (e *SourceReconnectRetryTemporaryError) Error() string { return fmt.Sprintf("temporary delivery for %s", e.Key) }
type SourceReconnectRetryStore struct { mu sync.Mutex; calls int; records []string }
func NewSourceReconnectRetryStore() *SourceReconnectRetryStore { return &SourceReconnectRetryStore{} }
func (s *SourceReconnectRetryStore) Attempt(key string) error {
    s.mu.Lock(); defer s.mu.Unlock(); s.calls++; s.records = append(s.records, key)
    if s.calls == 1 { return &SourceReconnectRetryTemporaryError{Key: key} }
    s.records = append(s.records, key); return nil
}
func (s *SourceReconnectRetryStore) Records() []string { s.mu.Lock(); defer s.mu.Unlock(); return append([]string(nil), s.records...) }
