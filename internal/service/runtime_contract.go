package service

import (
    "errors"
    "fmt"
    "log-aggregation/internal/store"
)

type SourceReconnectRetryCoordinator struct { backend *store.SourceReconnectRetryStore }
func NewSourceReconnectRetryCoordinator(b *store.SourceReconnectRetryStore) *SourceReconnectRetryCoordinator { return &SourceReconnectRetryCoordinator{backend: b} }
func (c *SourceReconnectRetryCoordinator) Send(key string) error {
    err := c.backend.Attempt(key)
    if err == nil { return nil }
    wrapped := fmt.Errorf("dispatch failed: %v", err)
    var temporary *store.SourceReconnectRetryTemporaryError
    if errors.As(wrapped, &temporary) { return c.backend.Attempt(key) }
    return wrapped
}
