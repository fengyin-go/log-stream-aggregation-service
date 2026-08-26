package service

import (
	"context"

	"log-aggregation/internal/store"
)

// BatchIngestCancellationCoordinator 按轮次向下游写入批量摄取的数据。
// 它必须尊重调用方的 ctx：一旦批量摄取超时或被取消，当前调用应尽快结束，
// 且不再安排下一轮下游写入。
type BatchIngestCancellationCoordinator struct {
	backend *store.BatchIngestCancellationStore
}

func NewBatchIngestCancellationCoordinator(b *store.BatchIngestCancellationStore) *BatchIngestCancellationCoordinator {
	return &BatchIngestCancellationCoordinator{backend: b}
}

func (c *BatchIngestCancellationCoordinator) Dispatch(ctx context.Context) error {
	for attempt := 0; attempt < 3; attempt++ {
		// 超时或取消后不再安排下一轮下游写入，当前调用尽快结束。
		if err := ctx.Err(); err != nil {
			return err
		}
		// 传入调用方的 ctx，使下游写入在超时后能立即返回，而非等完整延迟。
		if err := c.backend.Attempt(ctx); err != nil {
			return err
		}
	}
	return nil
}
