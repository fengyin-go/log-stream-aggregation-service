package service_test

import (
    "testing"
    "log-aggregation/internal/service"
    "log-aggregation/internal/store"
)

func TestR018Target(t *testing.T) {
    backend := store.NewRuleEvaluationAssemblyStore()
    coordinator := service.NewRuleEvaluationAssemblyCoordinator(backend)
    item, err := coordinator.Build("entry-1", true)
    cached, exists := backend.Get("entry-1")
    if err == nil || item != nil || exists || cached != nil {
        t.Fatalf("规则评估组装没有满足题面描述的触发顺序和公开预期")
    }
}
