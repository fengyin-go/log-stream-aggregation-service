package model

import (
	"time"
)

const (
	RetentionPolicyKeepAll  = "keep_all"
	RetentionPolicyKeep7d   = "keep_7d"
	RetentionPolicyKeep30d  = "keep_30d"
	RetentionPolicyKeep90d  = "keep_90d"
)

type RetentionPolicy struct {
	ID        string    `json:"id"`
	SourceID  string    `json:"source_id"`
	Policy    string    `json:"policy"`
	CreatedAt time.Time `json:"created_at"`
}

func (p *RetentionPolicy) Validate() error {
	if p.SourceID == "" {
		return NewValidationError("source_id", "日志源 ID 不能为空")
	}
	if p.Policy == "" {
		p.Policy = RetentionPolicyKeep30d
	}
	if p.Policy != RetentionPolicyKeepAll && p.Policy != RetentionPolicyKeep7d &&
		p.Policy != RetentionPolicyKeep30d && p.Policy != RetentionPolicyKeep90d {
		return NewValidationError("policy", "保留策略不合法")
	}
	return nil
}

type RetentionPolicyFilter struct {
	SourceID string
	Policy   string
}

func (f RetentionPolicyFilter) Match(p *RetentionPolicy) bool {
	if f.SourceID != "" && p.SourceID != f.SourceID {
		return false
	}
	if f.Policy != "" && p.Policy != f.Policy {
		return false
	}
	return true
}

type RetentionDays map[string]int

func GetRetentionDays(policy string) int {
	days := RetentionDays{
		RetentionPolicyKeepAll: 0,
		RetentionPolicyKeep7d:  7,
		RetentionPolicyKeep30d: 30,
		RetentionPolicyKeep90d: 90,
	}
	if d, ok := days[policy]; ok {
		return d
	}
	return 30
}
