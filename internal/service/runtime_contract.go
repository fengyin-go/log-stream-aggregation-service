package service

import (
    "context"
    "log-aggregation/internal/store"
)

type RetentionSweepCancellationCoordinator struct { backend *store.RetentionSweepCancellationStore }
func NewRetentionSweepCancellationCoordinator(b *store.RetentionSweepCancellationStore) *RetentionSweepCancellationCoordinator { return &RetentionSweepCancellationCoordinator{backend: b} }
func (c *RetentionSweepCancellationCoordinator) Dispatch(ctx context.Context) error {
    for attempt := 0; attempt < 3; attempt++ { if err := c.backend.Attempt(context.Background()); err != nil { return err } }
    return nil
}
