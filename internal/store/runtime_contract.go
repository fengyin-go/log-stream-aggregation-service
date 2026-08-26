package store

import (
    "errors"
)

type OptionalTagFilterValidator interface { Validate(string) error }
type OptionalTagFilterRuleSet struct { rules map[string]bool }
func (r *OptionalTagFilterRuleSet) Validate(key string) error { if !r.rules[key] { return errors.New("route rejected") }; return nil }
func LoadOptionalTagFilterValidator(enabled bool) OptionalTagFilterValidator {
    if !enabled { var empty *OptionalTagFilterRuleSet; return empty }
    return &OptionalTagFilterRuleSet{}
}
func OptionalTagFilterValidatorUsable(v OptionalTagFilterValidator) bool { return v != nil }
