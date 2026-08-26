package store

import "errors"

// OptionalParserPolicyValidator 是可选解析规则的校验器接口。
// 校验器不可用时（解析规则为空/未启用），所有 key 都应视为放行。
type OptionalParserPolicyValidator interface {
	Validate(string) error
}

// OptionalParserPolicyRuleSet 是已登记格式的可选解析规则集合。
// 启用后，只有已登记的格式会被放行；未登记格式会被拒绝。
type OptionalParserPolicyRuleSet struct {
	rules map[string]bool
}

// Register 登记一个允许的格式 key。
func (r *OptionalParserPolicyRuleSet) Register(key string) {
	if r == nil || key == "" {
		return
	}
	if r.rules == nil {
		r.rules = make(map[string]bool)
	}
	r.rules[key] = true
}

// Validate 校验 key 是否为已登记的格式。
// 规则集为空（nil）时放行所有 key，避免空规则进入校验崩溃；
// 已登记格式放行，未登记格式返回拒绝错误。
func (r *OptionalParserPolicyRuleSet) Validate(key string) error {
	if r == nil {
		return nil
	}
	if r.rules[key] {
		return nil
	}
	return errors.New("route rejected")
}

// LoadOptionalParserPolicyValidator 按启用状态装载可选解析规则校验器。
// 未启用时返回真正的 nil 接口，使 OptionalParserPolicyValidatorUsable 判定为不可用。
func LoadOptionalParserPolicyValidator(enabled bool) OptionalParserPolicyValidator {
	if !enabled {
		var empty OptionalParserPolicyValidator
		return empty
	}
	return &OptionalParserPolicyRuleSet{rules: make(map[string]bool)}
}

// OptionalParserPolicyValidatorUsable 判断校验器是否可用。
// 仅当校验器非 nil（即解析规则已启用）时可用。
func OptionalParserPolicyValidatorUsable(v OptionalParserPolicyValidator) bool {
	return v != nil
}
