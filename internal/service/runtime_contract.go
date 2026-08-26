package service

import (
    "context"
    "log-aggregation/internal/store"
)

type BatchIngestCancellationCoordinator struct { backend *store.BatchIngestCancellationStore }
func NewBatchIngestCancellationCoordinator(b *store.BatchIngestCancellationStore) *BatchIngestCancellationCoordinator { return &BatchIngestCancellationCoordinator{backend: b} }
func (c *BatchIngestCancellationCoordinator) Dispatch(ctx context.Context) error {
    for attempt := 0; attempt < 3; attempt++ { if err := c.backend.Attempt(context.Background()); err != nil { return err } }
    return nil
}
