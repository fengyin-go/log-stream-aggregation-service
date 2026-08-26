package service_test

import (
    "testing"
    "log-aggregation/internal/service"
    "log-aggregation/internal/store"
)

func TestR014Target(t *testing.T) {
    backend := store.NewNotificationDispatchRetryStore()
    coordinator := service.NewNotificationDispatchRetryCoordinator(backend)
    err := coordinator.Send("message-1")
    records := backend.Records()
    if err != nil || len(records) != 1 || records[0] != "message-1" {
        t.Fatalf("通知派发重试没有满足题面描述的触发顺序和公开预期")
    }
}
