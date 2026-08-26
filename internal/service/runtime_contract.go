package service

import (
    "context"
    "log-aggregation/internal/store"
)

type SourceHealthProbeCoordinator struct { backend *store.SourceHealthProbeStore }
func NewSourceHealthProbeCoordinator(b *store.SourceHealthProbeStore) *SourceHealthProbeCoordinator { return &SourceHealthProbeCoordinator{backend: b} }
func (c *SourceHealthProbeCoordinator) Probe(ctx context.Context, key string) error { return c.backend.Wait(context.Background(), key) }
