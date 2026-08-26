package store

import (
	"errors"

	"log-aggregation/internal/model"
)

// OptionalSeverityPolicyValidator 校验日志级别是否被可选级别规则登记。
type OptionalSeverityPolicyValidator interface{ Validate(string) error }

// OptionalSeverityPolicyRuleSet 持有已登记的合法日志级别集合。
type OptionalSeverityPolicyRuleSet struct{ rules map[string]bool }

func (r *OptionalSeverityPolicyRuleSet) Validate(key string) error {
	if r == nil || r.rules == nil || !r.rules[key] {
		return errors.New("route rejected")
	}
	return nil
}

// LoadOptionalSeverityPolicyValidator 按启用状态返回校验器。
//
// 关闭（enabled=false）时返回 nil 接口，使校验路径安全跳过，
// 避免返回“包裹 nil 指针的非 nil 接口”而在后续调用中解引用崩溃。
// 启用（enabled=true）时返回登记了全部合法日志级别的规则集，
// 使未登记的非法级别被拦截，而已登记级别正常放行。
func LoadOptionalSeverityPolicyValidator(enabled bool) OptionalSeverityPolicyValidator {
	if !enabled {
		return nil
	}
	rules := make(map[string]bool, len(model.ValidLogLevels()))
	for level := range model.ValidLogLevels() {
		rules[level] = true
	}
	return &OptionalSeverityPolicyRuleSet{rules: rules}
}

// OptionalSeverityPolicyValidatorUsable 判断校验器是否可用。
// nil 接口（关闭态）返回 false，包裹 nil 指针的接口也被视为不可用。
func OptionalSeverityPolicyValidatorUsable(v OptionalSeverityPolicyValidator) bool {
	if v == nil {
		return false
	}
	if rs, ok := v.(*OptionalSeverityPolicyRuleSet); ok && rs == nil {
		return false
	}
	return true
}
