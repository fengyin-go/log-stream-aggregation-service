package store

import (
	"fmt"
	"sync"
)

// AlertDeliveryRetryTemporaryError 表示告警投递过程中可重试的临时错误。
type AlertDeliveryRetryTemporaryError struct {
	Key string
}

func (e *AlertDeliveryRetryTemporaryError) Error() string {
	return fmt.Sprintf("temporary delivery for %s", e.Key)
}

// AlertDeliveryRetryStore 是告警投递后端的内存实现，用于验证重试协调器。
// 首次投递模拟临时失败，重试后成功；成功投递只落一条记录。
type AlertDeliveryRetryStore struct {
	mu      sync.Mutex
	calls   int
	records []string
}

func NewAlertDeliveryRetryStore() *AlertDeliveryRetryStore {
	return &AlertDeliveryRetryStore{}
}

// Attempt 模拟一次告警投递：首次返回临时错误且不落记录，重试成功后写入唯一一条记录。
func (s *AlertDeliveryRetryStore) Attempt(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls++
	if s.calls == 1 {
		// 首次投递模拟临时失败：不写入记录，避免失败后留下残留。
		return &AlertDeliveryRetryTemporaryError{Key: key}
	}
	// 重试成功：每个 key 只追加一条记录，杜绝重复。
	s.records = append(s.records, key)
	return nil
}

// Records 返回已成功投递的记录快照。
func (s *AlertDeliveryRetryStore) Records() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.records...)
}
