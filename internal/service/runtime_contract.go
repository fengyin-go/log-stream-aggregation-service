package service

import (
    "errors"
    "fmt"
    "log-aggregation/internal/store"
)

type NotificationDispatchRetryCoordinator struct { backend *store.NotificationDispatchRetryStore }
func NewNotificationDispatchRetryCoordinator(b *store.NotificationDispatchRetryStore) *NotificationDispatchRetryCoordinator { return &NotificationDispatchRetryCoordinator{backend: b} }
func (c *NotificationDispatchRetryCoordinator) Send(key string) error {
    err := c.backend.Attempt(key)
    if err == nil { return nil }
    var temporary *store.NotificationDispatchRetryTemporaryError
    if errors.As(err, &temporary) {
        // 临时故障需要重试，重试成功后底层仅保留一份记录。
        return c.backend.Attempt(key)
    }
    return fmt.Errorf("dispatch failed: %w", err)
}
