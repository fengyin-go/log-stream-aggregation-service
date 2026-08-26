package service

import "log-aggregation/internal/store"

// ValidateOptionalParserPolicy 校验可选解析规则。
// 校验器不可用（解析规则为空/未启用）时直接放行，避免空规则进入校验崩溃；
// 校验器可用时按已登记格式校验 key。
func ValidateOptionalParserPolicy(v store.OptionalParserPolicyValidator, key string) (err error) {
	if !store.OptionalParserPolicyValidatorUsable(v) {
		return nil
	}
	return v.Validate(key)
}
