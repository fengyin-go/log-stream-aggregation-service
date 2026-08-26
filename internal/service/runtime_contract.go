package service

import (
    "errors"
    "fmt"
    "log-aggregation/internal/store"
)

type AlertDeliveryRetryCoordinator struct { backend *store.AlertDeliveryRetryStore }
func NewAlertDeliveryRetryCoordinator(b *store.AlertDeliveryRetryStore) *AlertDeliveryRetryCoordinator { return &AlertDeliveryRetryCoordinator{backend: b} }
func (c *AlertDeliveryRetryCoordinator) Send(key string) error {
    err := c.backend.Attempt(key)
    if err == nil { return nil }
    wrapped := fmt.Errorf("dispatch failed: %v", err)
    var temporary *store.AlertDeliveryRetryTemporaryError
    if errors.As(wrapped, &temporary) { return c.backend.Attempt(key) }
    return wrapped
}
