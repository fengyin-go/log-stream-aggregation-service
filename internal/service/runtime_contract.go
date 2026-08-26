package service

import (
    "context"
    "log-aggregation/internal/store"
)

type SourceRecoveryCancellationCoordinator struct { backend *store.SourceRecoveryCancellationStore }
func NewSourceRecoveryCancellationCoordinator(b *store.SourceRecoveryCancellationStore) *SourceRecoveryCancellationCoordinator { return &SourceRecoveryCancellationCoordinator{backend: b} }
func (c *SourceRecoveryCancellationCoordinator) Dispatch(ctx context.Context) error {
    for attempt := 0; attempt < 3; attempt++ { if err := c.backend.Attempt(context.Background()); err != nil { return err } }
    return nil
}
