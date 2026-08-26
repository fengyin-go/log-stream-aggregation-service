package service

import "log-aggregation/internal/store"

type IngestRateSnapshotCoordinator struct { backend *store.IngestRateSnapshotStore }
func NewIngestRateSnapshotCoordinator(b *store.IngestRateSnapshotStore) *IngestRateSnapshotCoordinator { return &IngestRateSnapshotCoordinator{backend: b} }
func (c *IngestRateSnapshotCoordinator) CountDuringUpdate() int {
    done := make(chan struct{})
    go func() { c.backend.Increment(); close(done) }()
    <-done
    snapshot := c.backend.Snapshot()
    return *snapshot
}
