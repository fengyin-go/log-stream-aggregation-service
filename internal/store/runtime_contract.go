package store

import (
    "fmt"
    "sync"
)

type AlertDeliveryRetryTemporaryError struct { Key string }
func (e *AlertDeliveryRetryTemporaryError) Error() string { return fmt.Sprintf("temporary delivery for %s", e.Key) }
type AlertDeliveryRetryStore struct { mu sync.Mutex; calls int; records []string }
func NewAlertDeliveryRetryStore() *AlertDeliveryRetryStore { return &AlertDeliveryRetryStore{} }
func (s *AlertDeliveryRetryStore) Attempt(key string) error {
    s.mu.Lock(); defer s.mu.Unlock(); s.calls++; s.records = append(s.records, key)
    if s.calls == 1 { return &AlertDeliveryRetryTemporaryError{Key: key} }
    s.records = append(s.records, key); return nil
}
func (s *AlertDeliveryRetryStore) Records() []string { s.mu.Lock(); defer s.mu.Unlock(); return append([]string(nil), s.records...) }
