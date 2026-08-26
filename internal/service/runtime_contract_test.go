package service_test

import (
    "context"
    "errors"
    "testing"
    "time"
    "log-aggregation/internal/service"
    "log-aggregation/internal/store"
)

func TestR020Target(t *testing.T) {
    backend := store.NewBatchIngestCancellationStore(20 * time.Millisecond)
    coordinator := service.NewBatchIngestCancellationCoordinator(backend)
    stopped, stop := context.WithCancel(context.Background())
    stop()
    stoppedErr := coordinator.Dispatch(stopped)
    if !errors.Is(stoppedErr, context.Canceled) || backend.Calls() != 0 {
        t.Fatalf("请求在进入底层前已经取消却仍然开始了调用")
    }
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
    defer cancel()
    started := time.Now()
    err := coordinator.Dispatch(ctx)
    if !errors.Is(err, context.DeadlineExceeded) || time.Since(started) > 35*time.Millisecond || backend.Calls() != 1 {
        t.Fatalf("批量摄取取消没有满足题面描述的触发顺序和公开预期")
    }
}
