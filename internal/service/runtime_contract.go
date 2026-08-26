package service

import "log-aggregation/internal/store"

func ValidateOptionalParserPolicy(v store.OptionalParserPolicyValidator, key string) (err error) {
    
    if !store.OptionalParserPolicyValidatorUsable(v) { return nil }
    return v.Validate(key)
}
