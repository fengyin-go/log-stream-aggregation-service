package service_test

import (
    "testing"
    "log-aggregation/internal/service"
    "log-aggregation/internal/store"
)

func TestR025Target(t *testing.T) {
    pool := store.NewRetentionCursorLeasePool(2)
    coordinator := service.NewRetentionCursorLeaseCoordinator(pool)
    processed, err := coordinator.Process([]string{"a", "b", "c", "d"})
    if err != nil || processed != 4 || pool.Open() != 0 {
        t.Fatalf("清理游标租约没有满足题面描述的触发顺序和公开预期")
    }
}
