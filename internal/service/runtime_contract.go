package service

import (
    "context"
    "log-aggregation/internal/store"
)

type RetentionSweepCancellationCoordinator struct { backend *store.RetentionSweepCancellationStore }
func NewRetentionSweepCancellationCoordinator(b *store.RetentionSweepCancellationStore) *RetentionSweepCancellationCoordinator { return &RetentionSweepCancellationCoordinator{backend: b} }
func (c *RetentionSweepCancellationCoordinator) Dispatch(ctx context.Context) error {
    for attempt := 0; attempt < 3; attempt++ {
        // 请求超时或取消后立刻停住当前等待，剩余轮次不再继续，底层扫描调用次数不再上涨。
        if err := c.backend.Attempt(ctx); err != nil {
            return err
        }
    }
    return nil
}
