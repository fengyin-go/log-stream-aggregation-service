package service_test

import (
    "testing"
    "log-aggregation/internal/service"
    "log-aggregation/internal/store"
)

func callOptionalSeverityPolicy(v store.OptionalSeverityPolicyValidator, key string) (panicked bool, err error) {
    defer func() { if recover() != nil { panicked = true } }()
    err = service.ValidateOptionalSeverityPolicy(v, key)
    return
}

func TestR002Target(t *testing.T) {
    disabledPanic, disabledErr := callOptionalSeverityPolicy(store.LoadOptionalSeverityPolicyValidator(false), "anything")
    enabledPanic, enabledErr := callOptionalSeverityPolicy(store.LoadOptionalSeverityPolicyValidator(true), "missing")
    if disabledPanic || disabledErr != nil || enabledPanic || enabledErr == nil {
        t.Fatalf("可选级别规则没有满足题面描述的触发顺序和公开预期")
    }
}
