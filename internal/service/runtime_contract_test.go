package service_test

import (
    "context"
    "errors"
    "testing"
    "time"
    "log-aggregation/internal/service"
    "log-aggregation/internal/store"
)

func TestR001Target(t *testing.T) {
    backend := store.NewSourceHealthProbeStore(45 * time.Millisecond)
    coordinator := service.NewSourceHealthProbeCoordinator(backend)
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
    defer cancel()
    started := time.Now()
    firstErr := coordinator.Probe(ctx, "first")
    elapsed := time.Since(started)
    secondErr := coordinator.Probe(context.Background(), "second")
    if !errors.Is(firstErr, context.DeadlineExceeded) || elapsed > 30*time.Millisecond || secondErr != nil {
        t.Fatalf("来源健康探测没有满足题面描述的触发顺序和公开预期")
    }
}
