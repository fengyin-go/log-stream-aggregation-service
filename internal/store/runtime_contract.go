package store

import (
    "errors"
)

type OptionalTagFilterValidator interface { Validate(string) error }
type OptionalTagFilterRuleSet struct { rules map[string]bool }

// Validate 按 allowlist 校验 key：命中放行，未命中或规则集不可用则拒绝。
// fail-closed（nil 接收者或 nil rules 一律拒绝），既不 panic 也不漏检。
func (r *OptionalTagFilterRuleSet) Validate(key string) error {
    if r == nil || r.rules == nil || !r.rules[key] {
        return errors.New("route rejected")
    }
    return nil
}

// LoadOptionalTagFilterValidator 关闭时返回 untyped nil 接口（而非 typed-nil 指针），
// 使 OptionalTagFilterValidatorUsable 正确判定为不可用，调用方安全跳过；
// 开启时返回 rules 已初始化的规则集，支持后续填充放行标签。
func LoadOptionalTagFilterValidator(enabled bool) OptionalTagFilterValidator {
    if !enabled {
        return nil
    }
    return &OptionalTagFilterRuleSet{rules: make(map[string]bool)}
}

func OptionalTagFilterValidatorUsable(v OptionalTagFilterValidator) bool { return v != nil }
