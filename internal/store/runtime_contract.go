package store

import (
    "errors"
)

type OptionalSeverityPolicyValidator interface { Validate(string) error }
type OptionalSeverityPolicyRuleSet struct { rules map[string]bool }
func (r *OptionalSeverityPolicyRuleSet) Validate(key string) error { if !r.rules[key] { return errors.New("route rejected") }; return nil }
func LoadOptionalSeverityPolicyValidator(enabled bool) OptionalSeverityPolicyValidator {
    if !enabled { var empty *OptionalSeverityPolicyRuleSet; return empty }
    return &OptionalSeverityPolicyRuleSet{}
}
func OptionalSeverityPolicyValidatorUsable(v OptionalSeverityPolicyValidator) bool { return v != nil }
