package store

import (
    "errors"
)

type OptionalParserPolicyValidator interface { Validate(string) error }
type OptionalParserPolicyRuleSet struct { rules map[string]bool }
func (r *OptionalParserPolicyRuleSet) Validate(key string) error { if !r.rules[key] { return errors.New("route rejected") }; return nil }
func LoadOptionalParserPolicyValidator(enabled bool) OptionalParserPolicyValidator {
    if !enabled { var empty *OptionalParserPolicyRuleSet; return empty }
    return &OptionalParserPolicyRuleSet{}
}
func OptionalParserPolicyValidatorUsable(v OptionalParserPolicyValidator) bool { return v != nil }
