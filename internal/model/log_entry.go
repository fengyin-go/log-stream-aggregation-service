package model

import (
	"strings"
	"time"
)

const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
	LogLevelFatal = "fatal"
)

var validLogLevels = map[string]bool{
	LogLevelDebug: true,
	LogLevelInfo:  true,
	LogLevelWarn:  true,
	LogLevelError: true,
	LogLevelFatal: true,
}

type LogEntry struct {
	ID        string    `json:"id"`
	SourceID  string    `json:"source_id"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
}

func (e *LogEntry) Validate() error {
	e.Message = strings.TrimSpace(e.Message)
	if e.SourceID == "" {
		return NewValidationError("source_id", "日志源 ID 不能为空")
	}
	if e.Message == "" {
		return NewValidationError("message", "日志内容不能为空")
	}
	if e.Level == "" {
		e.Level = LogLevelInfo
	}
	if !validLogLevels[e.Level] {
		return NewValidationError("level", "日志级别不合法")
	}
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}
	return nil
}

type LogEntryFilter struct {
	SourceID string
	Level    string
	Keyword  string
	Tag      string
	StartAt  time.Time
	EndAt    time.Time
}

func (f LogEntryFilter) Match(e *LogEntry) bool {
	if f.SourceID != "" && e.SourceID != f.SourceID {
		return false
	}
	if f.Level != "" && e.Level != f.Level {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(e.Message), k) {
			return false
		}
	}
	if f.Tag != "" {
		found := false
		for _, t := range e.Tags {
			if t == f.Tag {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	if !f.StartAt.IsZero() && e.Timestamp.Before(f.StartAt) {
		return false
	}
	if !f.EndAt.IsZero() && e.Timestamp.After(f.EndAt) {
		return false
	}
	return true
}
