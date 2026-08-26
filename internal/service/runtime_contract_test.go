package service_test

import (
    "testing"
    "log-aggregation/internal/service"
    "log-aggregation/internal/store"
)

func callOptionalParserPolicy(v store.OptionalParserPolicyValidator, key string) (panicked bool, err error) {
    defer func() { if recover() != nil { panicked = true } }()
    err = service.ValidateOptionalParserPolicy(v, key)
    return
}

func TestR022Target(t *testing.T) {
    disabledPanic, disabledErr := callOptionalParserPolicy(store.LoadOptionalParserPolicyValidator(false), "anything")
    enabledPanic, enabledErr := callOptionalParserPolicy(store.LoadOptionalParserPolicyValidator(true), "missing")
    if disabledPanic || disabledErr != nil || enabledPanic || enabledErr == nil {
        t.Fatalf("可选解析规则没有满足题面描述的触发顺序和公开预期")
    }
}
