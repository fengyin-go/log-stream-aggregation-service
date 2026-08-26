package service_test

import (
    "testing"
    "time"
    "log-aggregation/internal/service"
    "log-aggregation/internal/store"
)

func TestR006Target(t *testing.T) {
    coordinator := service.NewSourceScanResultStreamCoordinator(store.NewSourceScanResultStreamStore())
    done := make(chan error, 1)
    go func() { _, err := coordinator.Collect(true); done <- err }()
    select {
    case err := <-done:
        if err == nil { t.Fatalf("分区失败没有返回错误并结束结果流") }
    case <-time.After(80 * time.Millisecond):
        t.Fatalf("来源扫描结果流没有满足题面描述的触发顺序和公开预期")
    }
}
