package service_test

import (
    "testing"
    "log-aggregation/internal/service"
    "log-aggregation/internal/store"
)

func TestR029Target(t *testing.T) {
    coordinator := service.NewIngestRateSnapshotCoordinator(store.NewIngestRateSnapshotStore(1))
    if coordinator.CountDuringUpdate() != 1 {
        t.Fatalf("摄取速率快照没有满足题面描述的触发顺序和公开预期")
    }
}
