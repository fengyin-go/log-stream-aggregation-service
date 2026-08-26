package service_test

import (
    "testing"
    "log-aggregation/internal/service"
    "log-aggregation/internal/store"
)

func TestR009Target(t *testing.T) {
    coordinator := service.NewSourceCounterSnapshotCoordinator(store.NewSourceCounterSnapshotStore(1))
    if coordinator.CountDuringUpdate() != 1 {
        t.Fatalf("来源计数快照没有满足题面描述的触发顺序和公开预期")
    }
}
