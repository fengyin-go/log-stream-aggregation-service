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
    s.mu.Lock(); defer s.mu.Unlock(); s.calls++
    if s.calls == 1 {
        // 首次派发返回临时故障，此时尚未真正投递，不应留下任何记录，
        // 避免底层残留不完整记录并在重试时产生重复。
        return &NotificationDispatchRetryTemporaryError{Key: key}
    }
    // 重试成功后才记入一条投递记录，保证成功后只保留一份。
    s.records = append(s.records, key)
    return nil
}
func (s *NotificationDispatchRetryStore) Records() []string { s.mu.Lock(); defer s.mu.Unlock(); return append([]string(nil), s.records...) }
