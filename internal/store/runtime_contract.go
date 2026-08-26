package store

import (
    "context"
    "sync"
    "time"
)

type RetentionSweepCancellationStore struct { mu sync.Mutex; calls int; delay time.Duration }
func NewRetentionSweepCancellationStore(delay time.Duration) *RetentionSweepCancellationStore { return &RetentionSweepCancellationStore{delay: delay} }
func (s *RetentionSweepCancellationStore) Attempt(ctx context.Context) error {
    // 请求开始前就已取消：底层扫描一次都不调用，调用次数不上涨。
    if err := ctx.Err(); err != nil {
        return err
    }
    s.mu.Lock(); s.calls++; s.mu.Unlock()
    timer := time.NewTimer(s.delay); defer timer.Stop()
    // 已开始的请求超时后立刻停住当前等待，后续轮次不再继续，调用次数不再上涨。
    select { case <-ctx.Done(): return ctx.Err(); case <-timer.C: return nil }
}
func (s *RetentionSweepCancellationStore) Calls() int { s.mu.Lock(); defer s.mu.Unlock(); return s.calls }
