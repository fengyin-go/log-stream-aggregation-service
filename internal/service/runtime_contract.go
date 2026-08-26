package service

import (
	"context"
	"log-aggregation/internal/store"
)

type SourceRecoveryCancellationCoordinator struct {
	backend *store.SourceRecoveryCancellationStore
}

func NewSourceRecoveryCancellationCoordinator(b *store.SourceRecoveryCancellationStore) *SourceRecoveryCancellationCoordinator {
	return &SourceRecoveryCancellationCoordinator{backend: b}
}
func (c *SourceRecoveryCancellationCoordinator) Dispatch(ctx context.Context) error {
	for attempt := 0; attempt < 3; attempt++ {
		// 取消信号在启动新一轮次前先行短路，避免后续恢复轮次继续启动。
		if err := ctx.Err(); err != nil {
			return err
		}
		// 将 ctx 透传到底层，使正在进行的轮次能在取消时即时中断。
		if err := c.backend.Attempt(ctx); err != nil {
			return err
		}
	}
	return nil
}
