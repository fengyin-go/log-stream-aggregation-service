package service_test

import (
    "testing"
    "log-aggregation/internal/service"
    "log-aggregation/internal/store"
)

func TestR003Target(t *testing.T) {
    backend := store.NewIngestPayloadSnapshotStore()
    coordinator := service.NewIngestPayloadSnapshotCoordinator(backend)
    payload := []byte("alpha")
    coordinator.Archive(payload)
    copy(payload, []byte("omega"))
    if string(backend.Snapshot()) != "alpha" || string(coordinator.Export()) != "alpha" {
        t.Fatalf("摄取内容快照没有满足题面描述的触发顺序和公开预期")
    }
}
