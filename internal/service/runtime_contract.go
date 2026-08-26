package service

import "log-aggregation/internal/store"

func ValidateOptionalSeverityPolicy(v store.OptionalSeverityPolicyValidator, key string) (err error) {
    
    if !store.OptionalSeverityPolicyValidatorUsable(v) { return nil }
    return v.Validate(key)
}
