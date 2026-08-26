package store

import (
	"context"
	"time"
)

// TailCursorProbeStore 对尾部游标进行长轮询探测：在延迟窗口内等待游标推进，
// 到期则返回 nil 表示暂无新数据。每次 Wait 的等待状态彼此完全独立——
// 仅监听本次调用传入的 ctx，不缓存或复用任何先前请求的上下文，
// 避免一次探测的截止时间泄漏到后续的无超时探测上。
type TailCursorProbeStore struct {
	delay time.Duration
}

func NewTailCursorProbeStore(delay time.Duration) *TailCursorProbeStore {
	return &TailCursorProbeStore{delay: delay}
}

func (s *TailCursorProbeStore) Wait(ctx context.Context, key string) error {
	timer := time.NewTimer(s.delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
