package service_test

import (
    "testing"
    "log-aggregation/internal/service"
    "log-aggregation/internal/store"
)

func TestR013Target(t *testing.T) {
    backend := store.NewQueryTokenArchiveStore()
    coordinator := service.NewQueryTokenArchiveCoordinator(backend)
    payload := []byte("alpha")
    coordinator.Archive(payload)
    copy(payload, []byte("omega"))
    if string(backend.Snapshot()) != "alpha" || string(coordinator.Export()) != "alpha" {
        t.Fatalf("查询令牌归档没有满足题面描述的触发顺序和公开预期")
    }
}
