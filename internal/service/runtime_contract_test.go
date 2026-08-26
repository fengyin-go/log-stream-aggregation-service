package service_test

import (
    "testing"
    "log-aggregation/internal/service"
    "log-aggregation/internal/store"
)

func TestR028Target(t *testing.T) {
    backend := store.NewLogMetadataAssemblyStore()
    coordinator := service.NewLogMetadataAssemblyCoordinator(backend)
    item, err := coordinator.Build("entry-1", true)
    cached, exists := backend.Get("entry-1")
    if err == nil || item != nil || exists || cached != nil {
        t.Fatalf("日志元数据组装没有满足题面描述的触发顺序和公开预期")
    }
}
