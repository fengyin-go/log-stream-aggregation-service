package store

import (
    "fmt"
    "sync"
)

type NotificationDispatchRetryTemporaryError struct { Key string }
func (e *NotificationDispatchRetryTemporaryError) Error() string { return fmt.Sprintf("temporary delivery for %s", e.Key) }
type NotificationDispatchRetryStore struct { mu sync.Mutex; calls int; records []string }
func NewNotificationDispatchRetryStore() *NotificationDispatchRetryStore { return &NotificationDispatchRetryStore{} }
func (s *NotificationDispatchRetryStore) Attempt(key string) error {
    s.mu.Lock(); defer s.mu.Unlock(); s.calls++; s.records = append(s.records, key)
    if s.calls == 1 { return &NotificationDispatchRetryTemporaryError{Key: key} }
    s.records = append(s.records, key); return nil
}
func (s *NotificationDispatchRetryStore) Records() []string { s.mu.Lock(); defer s.mu.Unlock(); return append([]string(nil), s.records...) }
