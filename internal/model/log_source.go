package model

import (
	"strings"
	"time"
)

const (
	LogSourceTypeApp    = "app"
	LogSourceTypeService = "service"
	LogSourceTypeDB     = "db"

	LogSourceStatusActive  = "active"
	LogSourceStatusPaused = "paused"
)

type LogSource struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Host      string    `json:"host"`
	Path      string    `json:"path"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *LogSource) Validate() error {
	s.Name = strings.TrimSpace(s.Name)
	s.Host = strings.TrimSpace(s.Host)
	s.Path = strings.TrimSpace(s.Path)
	if s.Name == "" {
		return NewValidationError("name", "日志源名称不能为空")
	}
	if s.Type == "" {
		s.Type = LogSourceTypeApp
	}
	if s.Type != LogSourceTypeApp && s.Type != LogSourceTypeService && s.Type != LogSourceTypeDB {
		return NewValidationError("type", "日志源类型不合法")
	}
	if s.Host == "" {
		return NewValidationError("host", "主机地址不能为空")
	}
	if s.Path == "" {
		return NewValidationError("path", "日志路径不能为空")
	}
	if s.Status == "" {
		s.Status = LogSourceStatusActive
	}
	if s.Status != LogSourceStatusActive && s.Status != LogSourceStatusPaused {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

func LogSourceCanTransition(from, to string) bool {
	if from == LogSourceStatusActive && to == LogSourceStatusPaused {
		return true
	}
	if from == LogSourceStatusPaused && to == LogSourceStatusActive {
		return true
	}
	return false
}

type LogSourceFilter struct {
	Name     string
	Type     string
	Status   string
	Keyword  string
}

func (f LogSourceFilter) Match(s *LogSource) bool {
	if f.Name != "" && s.Name != f.Name {
		return false
	}
	if f.Type != "" && s.Type != f.Type {
		return false
	}
	if f.Status != "" && s.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(s.Name), k) &&
			!strings.Contains(strings.ToLower(s.Host), k) {
			return false
		}
	}
	return true
}
