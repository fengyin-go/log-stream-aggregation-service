package service

import (
    "context"
    "log-aggregation/internal/store"
)

type RetentionPolicyProbeCoordinator struct { backend *store.RetentionPolicyProbeStore }
func NewRetentionPolicyProbeCoordinator(b *store.RetentionPolicyProbeStore) *RetentionPolicyProbeCoordinator { return &RetentionPolicyProbeCoordinator{backend: b} }
func (c *RetentionPolicyProbeCoordinator) Probe(ctx context.Context, key string) error { return c.backend.Wait(ctx, key) }
