package service_test

import (
    "testing"
    "log-aggregation/internal/service"
    "log-aggregation/internal/store"
)

func TestR017Target(t *testing.T) {
    backend := store.NewSourcePauseVersionStore()
    coordinator := service.NewSourcePauseVersionCoordinator(backend)
    coordinator.CompleteThenLate("delivery-1")
    state := backend.Get("delivery-1")
    if state.State != "complete" || state.Version != 2 || coordinator.EffectCount() != 1 {
        t.Fatalf("来源暂停终态没有满足题面描述的触发顺序和公开预期")
    }
}
