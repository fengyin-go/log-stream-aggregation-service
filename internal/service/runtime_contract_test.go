package service_test

import (
    "testing"
    "log-aggregation/internal/service"
    "log-aggregation/internal/store"
)

func callOptionalTagFilter(v store.OptionalTagFilterValidator, key string) (panicked bool, err error) {
    defer func() { if recover() != nil { panicked = true } }()
    err = service.ValidateOptionalTagFilter(v, key)
    return
}

func TestR012Target(t *testing.T) {
    disabledPanic, disabledErr := callOptionalTagFilter(store.LoadOptionalTagFilterValidator(false), "anything")
    enabledPanic, enabledErr := callOptionalTagFilter(store.LoadOptionalTagFilterValidator(true), "missing")
    if disabledPanic || disabledErr != nil || enabledPanic || enabledErr == nil {
        t.Fatalf("可选标签过滤没有满足题面描述的触发顺序和公开预期")
    }
}
