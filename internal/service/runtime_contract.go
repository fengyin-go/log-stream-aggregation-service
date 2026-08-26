package service

import "log-aggregation/internal/store"

type ErrorRateSnapshotCoordinator struct { backend *store.ErrorRateSnapshotStore }
func NewErrorRateSnapshotCoordinator(b *store.ErrorRateSnapshotStore) *ErrorRateSnapshotCoordinator { return &ErrorRateSnapshotCoordinator{backend: b} }
func (c *ErrorRateSnapshotCoordinator) CountDuringUpdate() int {
    done := make(chan struct{})
    go func() { c.backend.Increment(); close(done) }()
    <-done
    snapshot := c.backend.Snapshot()
    return *snapshot
}
