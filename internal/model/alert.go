package model

import (
	"strings"
	"time"
)

const (
	AlertStatusOpen         = "open"
	AlertStatusAcknowledged = "acknowledged"
	AlertStatusResolved     = "resolved"
)

var alertTransitions = map[string]map[string]bool{
	AlertStatusOpen:         {AlertStatusAcknowledged: true, AlertStatusResolved: true},
	AlertStatusAcknowledged: {AlertStatusResolved: true},
}

func AlertCanTransition(from, to string) bool {
	if m, ok := alertTransitions[from]; ok {
		return m[to]
	}
	return false
}

type Alert struct {
	ID        string    `json:"id"`
	RuleID    string    `json:"rule_id"`
	SourceID  string    `json:"source_id"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (a *Alert) Validate() error {
	a.Message = strings.TrimSpace(a.Message)
	if a.RuleID == "" {
		return NewValidationError("rule_id", "规则 ID 不能为空")
	}
	if a.SourceID == "" {
		return NewValidationError("source_id", "日志源 ID 不能为空")
	}
	if a.Message == "" {
		return NewValidationError("message", "告警消息不能为空")
	}
	if a.Status == "" {
		a.Status = AlertStatusOpen
	}
	if a.Status != AlertStatusOpen && a.Status != AlertStatusAcknowledged && a.Status != AlertStatusResolved {
		return NewValidationError("status", "告警状态不合法")
	}
	if a.Level == "" {
		a.Level = LogLevelError
	}
	if !validLogLevels[a.Level] {
		return NewValidationError("level", "告警级别不合法")
	}
	return nil
}

type AlertFilter struct {
	RuleID   string
	SourceID string
	Level    string
	Status   string
	Keyword  string
}

func (f AlertFilter) Match(a *Alert) bool {
	if f.RuleID != "" && a.RuleID != f.RuleID {
		return false
	}
	if f.SourceID != "" && a.SourceID != f.SourceID {
		return false
	}
	if f.Level != "" && a.Level != f.Level {
		return false
	}
	if f.Status != "" && a.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(a.Message), k) {
			return false
		}
	}
	return true
}
