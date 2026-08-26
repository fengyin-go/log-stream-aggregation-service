package model

import (
	"strings"
	"time"
)

const (
	AlertRuleStatusActive  = "active"
	AlertRuleStatusPaused = "paused"
)

type AlertRule struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	SourceID       string    `json:"source_id"`
	LevelThreshold string    `json:"level_threshold"`
	Keyword        string    `json:"keyword"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

func (r *AlertRule) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	if r.Name == "" {
		return NewValidationError("name", "规则名称不能为空")
	}
	if r.SourceID == "" {
		return NewValidationError("source_id", "日志源 ID 不能为空")
	}
	if r.LevelThreshold == "" {
		return NewValidationError("level_threshold", "告警级别阈值不能为空")
	}
	if !validLogLevels[r.LevelThreshold] {
		return NewValidationError("level_threshold", "级别阈值不合法")
	}
	if r.Status == "" {
		r.Status = AlertRuleStatusActive
	}
	if r.Status != AlertRuleStatusActive && r.Status != AlertRuleStatusPaused {
		return NewValidationError("status", "规则状态不合法")
	}
	return nil
}

func AlertRuleCanTransition(from, to string) bool {
	if from == AlertRuleStatusActive && to == AlertRuleStatusPaused {
		return true
	}
	if from == AlertRuleStatusPaused && to == AlertRuleStatusActive {
		return true
	}
	return false
}

type AlertRuleFilter struct {
	Name     string
	SourceID string
	Status   string
	Keyword  string
}

func (f AlertRuleFilter) Match(r *AlertRule) bool {
	if f.Name != "" && r.Name != f.Name {
		return false
	}
	if f.SourceID != "" && r.SourceID != f.SourceID {
		return false
	}
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(r.Name), k) &&
			!strings.Contains(strings.ToLower(r.Keyword), k) {
			return false
		}
	}
	return true
}
