package service_test

import (
    "testing"
    "log-aggregation/internal/service"
    "log-aggregation/internal/store"
)

func TestR015Target(t *testing.T) {
    pool := store.NewIndexReaderLeasePool(2)
    coordinator := service.NewIndexReaderLeaseCoordinator(pool)
    processed, err := coordinator.Process([]string{"a", "b", "c", "d"})
    if err != nil || processed != 4 || pool.Open() != 0 {
        t.Fatalf("索引读取租约没有满足题面描述的触发顺序和公开预期")
    }
}
