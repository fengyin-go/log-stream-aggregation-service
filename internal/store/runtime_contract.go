package store

import (
    "context"
    "time"
)

type RetentionPolicyProbeStore struct { delay time.Duration }
func NewRetentionPolicyProbeStore(delay time.Duration) *RetentionPolicyProbeStore { return &RetentionPolicyProbeStore{delay: delay} }

// Wait 等待固定的探测延迟。每次探测只受自身 ctx 的超时/取消控制，
// 不复用其它探测的上下文，避免某个探测取消后影响后续探测的返回时机。
func (s *RetentionPolicyProbeStore) Wait(ctx context.Context, key string) error {
    timer := time.NewTimer(s.delay)
    defer timer.Stop()
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-timer.C:
        return nil
    }
}
