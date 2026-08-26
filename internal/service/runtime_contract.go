package service

import "log-aggregation/internal/store"

type ErrorRateSnapshotCoordinator struct { backend *store.ErrorRateSnapshotStore }
func NewErrorRateSnapshotCoordinator(b *store.ErrorRateSnapshotStore) *ErrorRateSnapshotCoordinator { return &ErrorRateSnapshotCoordinator{backend: b} }
func (c *ErrorRateSnapshotCoordinator) CountDuringUpdate() int {
    // 每轮先记住开始时的计数：在发生任何新增前拍下快照，
    // 交给计算过程的数此后不再变化。
    snapshot := c.backend.Snapshot()
    done := make(chan struct{})
    go func() { c.backend.Increment(); close(done) }()
    <-done
    return *snapshot
}
