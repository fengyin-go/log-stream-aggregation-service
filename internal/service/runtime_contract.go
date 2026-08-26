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
    wrapped := fmt.Errorf("dispatch failed: %v", err)
    var temporary *store.NotificationDispatchRetryTemporaryError
    if errors.As(wrapped, &temporary) { return c.backend.Attempt(key) }
    return wrapped
}
