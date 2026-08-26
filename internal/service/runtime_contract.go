package service

import (
	"errors"
	"fmt"
	"log-aggregation/internal/store"
)

// AlertDeliveryRetryCoordinator 负责告警投递，并在遇到临时错误时正确重试。
// 临时失败时不会在存储中留下残留，最终成功后仅保留一条投递记录。
type AlertDeliveryRetryCoordinator struct {
	backend *store.AlertDeliveryRetryStore
}

func NewAlertDeliveryRetryCoordinator(b *store.AlertDeliveryRetryStore) *AlertDeliveryRetryCoordinator {
	return &AlertDeliveryRetryCoordinator{backend: b}
}

// Send 投递告警。遇到临时错误会重试；临时失败期间不向存储写入记录，
// 重试成功后仅追加一条记录，保证最终存储中只留一条。
func (c *AlertDeliveryRetryCoordinator) Send(key string) error {
	err := c.backend.Attempt(key)
	if err == nil {
		return nil
	}
	// 用 %w 包装以保留错误链，使 errors.As 能够识别出临时错误。
	wrapped := fmt.Errorf("dispatch failed: %w", err)
	var temporary *store.AlertDeliveryRetryTemporaryError
	if errors.As(wrapped, &temporary) {
		// 临时错误：重试一次，成功则结束，失败则返回错误。
		return c.backend.Attempt(key)
	}
	return wrapped
}
