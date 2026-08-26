package store

import (
	"context"
	"time"
)

// SourceHealthProbeStore 模拟来源健康探测的后台延迟。
type SourceHealthProbeStore struct{ delay time.Duration }

func NewSourceHealthProbeStore(delay time.Duration) *SourceHealthProbeStore {
	return &SourceHealthProbeStore{delay: delay}
}

// Wait 阻塞最多 delay 后返回 nil；ctx 超时或取消时立即返回其错误。
// 每次调用相互独立：一次探测的状态不会泄漏到后续探测。
func (s *SourceHealthProbeStore) Wait(ctx context.Context, key string) error {
	timer := time.NewTimer(s.delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
