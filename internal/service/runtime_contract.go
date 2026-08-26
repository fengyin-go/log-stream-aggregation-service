package service

import "log-aggregation/internal/store"

type SourceCounterSnapshotCoordinator struct { backend *store.SourceCounterSnapshotStore }
func NewSourceCounterSnapshotCoordinator(b *store.SourceCounterSnapshotStore) *SourceCounterSnapshotCoordinator { return &SourceCounterSnapshotCoordinator{backend: b} }
func (c *SourceCounterSnapshotCoordinator) CountDuringUpdate() int {
    done := make(chan struct{})
    go func() { c.backend.Increment(); close(done) }()
    <-done
    snapshot := c.backend.Snapshot()
    return *snapshot
}
