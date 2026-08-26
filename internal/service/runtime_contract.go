package service

import (
	"context"

	"log-aggregation/internal/store"
)

// SourceHealthProbeCoordinator 协调来源健康探测，尊重调用方传入的超时与取消。
type SourceHealthProbeCoordinator struct{ backend *store.SourceHealthProbeStore }

func NewSourceHealthProbeCoordinator(b *store.SourceHealthProbeStore) *SourceHealthProbeCoordinator {
	return &SourceHealthProbeCoordinator{backend: b}
}

// Probe 透传调用方 ctx，使超时与取消能在后台延迟结束前生效。
func (c *SourceHealthProbeCoordinator) Probe(ctx context.Context, key string) error {
	return c.backend.Wait(ctx, key)
}
