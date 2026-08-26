package service_test

import (
    "testing"
    "log-aggregation/internal/service"
    "log-aggregation/internal/store"
)

func TestR027Target(t *testing.T) {
    backend := store.NewNotificationReadVersionStore()
    coordinator := service.NewNotificationReadVersionCoordinator(backend)
    coordinator.CompleteThenLate("delivery-1")
    state := backend.Get("delivery-1")
    if state.State != "complete" || state.Version != 2 || coordinator.EffectCount() != 1 {
        t.Fatalf("通知已读终态没有满足题面描述的触发顺序和公开预期")
    }
}
