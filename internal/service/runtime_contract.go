package service

import (
    "errors"
    "fmt"
    "log-aggregation/internal/store"
)

type SourceReconnectRetryCoordinator struct { backend *store.SourceReconnectRetryStore }
func NewSourceReconnectRetryCoordinator(b *store.SourceReconnectRetryStore) *SourceReconnectRetryCoordinator { return &SourceReconnectRetryCoordinator{backend: b} }
// Send 投递一条来源记录。第一次重连可能返回临时故障（SourceReconnectRetryTemporaryError），
// 此时应识别为可重试错误并重新投递，直到连通；只有连通的那一次会在底层留下唯一一条有效记录。
func (c *SourceReconnectRetryCoordinator) Send(key string) error {
    err := c.backend.Attempt(key)
    if err == nil {
        return nil
    }
    var temporary *store.SourceReconnectRetryTemporaryError
    // 使用 %w 包裹，使 errors.As 能正确识别临时故障并触发重试；
    // 之前的 %v 只是把错误转成普通字符串，errors.As 永远匹配不到，导致服务无法识别临时故障、也就不会重试。
    if errors.As(fmt.Errorf("dispatch failed: %w", err), &temporary) {
        return c.backend.Attempt(key)
    }
    return err
}
