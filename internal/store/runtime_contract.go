package store

import (
    "fmt"
    "sync"
)

type SourceReconnectRetryTemporaryError struct { Key string }
func (e *SourceReconnectRetryTemporaryError) Error() string { return fmt.Sprintf("temporary delivery for %s", e.Key) }
type SourceReconnectRetryStore struct { mu sync.Mutex; calls int; records []string }
func NewSourceReconnectRetryStore() *SourceReconnectRetryStore { return &SourceReconnectRetryStore{} }
// SourceReconnectRetryStore 模拟来源重连重试的底层投递：第一次重连返回临时
// 故障，重试后连通。只有真正连通的那一次才落盘一条有效记录——返回临时故障的
// 尝试不写入任何记录，避免底层留下部分或重复记录。
func (s *SourceReconnectRetryStore) Attempt(key string) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.calls++
    if s.calls == 1 {
        return &SourceReconnectRetryTemporaryError{Key: key}
    }
    s.records = append(s.records, key)
    return nil
}
func (s *SourceReconnectRetryStore) Records() []string { s.mu.Lock(); defer s.mu.Unlock(); return append([]string(nil), s.records...) }
