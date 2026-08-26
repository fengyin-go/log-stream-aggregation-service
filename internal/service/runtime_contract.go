package service

import (
	"context"
	"log-aggregation/internal/store"
)

type TailCursorProbeCoordinator struct{ backend *store.TailCursorProbeStore }

func NewTailCursorProbeCoordinator(b *store.TailCursorProbeStore) *TailCursorProbeCoordinator {
	return &TailCursorProbeCoordinator{backend: b}
}

func (c *TailCursorProbeCoordinator) Probe(ctx context.Context, key string) error {
	return c.backend.Wait(ctx, key)
}
