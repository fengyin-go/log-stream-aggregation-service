package service

import "log-aggregation/internal/store"

func ValidateOptionalTagFilter(v store.OptionalTagFilterValidator, key string) (err error) {
    
    if !store.OptionalTagFilterValidatorUsable(v) { return nil }
    return v.Validate(key)
}
