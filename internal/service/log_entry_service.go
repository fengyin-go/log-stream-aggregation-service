package service

import (
	"sort"
	"strings"
	"time"

	"log-aggregation/internal/model"
	"log-aggregation/pkg/idgen"
)

func (s *Service) CreateLogEntry(input model.LogEntry) (*model.LogEntry, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	src, err := s.store.GetLogSource(input.SourceID)
	if err != nil {
		return nil, model.NewValidationError("source_id", "日志源不存在")
	}
	if src.Status == model.LogSourceStatusPaused {
		return nil, model.NewValidationError("source_id", "日志源已暂停，无法摄取")
	}
	entry := &model.LogEntry{
		ID:        idgen.Hex(),
		SourceID:  input.SourceID,
		Level:     input.Level,
		Message:   input.Message,
		Timestamp: input.Timestamp,
		Tags:      input.Tags,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateLogEntry(entry); err != nil {
		return nil, err
	}
	s.checkAlertRules(entry, src)
	s.log.Infof("创建日志条目: %s", entry.ID)
	return entry, nil
}

func (s *Service) BatchCreateLogEntries(inputs []model.LogEntry) ([]*model.LogEntry, error) {
	if len(inputs) == 0 {
		return nil, model.NewValidationError("entries", "批量日志不能为空")
	}
	entries := make([]*model.LogEntry, 0, len(inputs))
	for _, input := range inputs {
		if err := input.Validate(); err != nil {
			return nil, err
		}
		src, err := s.store.GetLogSource(input.SourceID)
		if err != nil {
			return nil, model.NewValidationError("source_id", "日志源不存在: "+input.SourceID)
		}
		if src.Status == model.LogSourceStatusPaused {
			return nil, model.NewValidationError("source_id", "日志源已暂停，无法摄取: "+input.SourceID)
		}
		entry := &model.LogEntry{
			ID:        idgen.Hex(),
			SourceID:  input.SourceID,
			Level:     input.Level,
			Message:   input.Message,
			Timestamp: input.Timestamp,
			Tags:      input.Tags,
			CreatedAt: time.Now(),
		}
		entries = append(entries, entry)
	}
	if err := s.store.BatchCreateLogEntries(entries); err != nil {
		return nil, err
	}
	for _, entry := range entries {
		src, _ := s.store.GetLogSource(entry.SourceID)
		if src != nil {
			s.checkAlertRules(entry, src)
		}
	}
	s.log.Infof("批量创建日志条目: %d 条", len(entries))
	return entries, nil
}

func (s *Service) checkAlertRules(entry *model.LogEntry, src *model.LogSource) {
	rules := s.store.ListActiveAlertRulesBySourceID(entry.SourceID)
	for _, rule := range rules {
		matched := false
		if rule.LevelThreshold != "" && levelRank(entry.Level) >= levelRank(rule.LevelThreshold) {
			matched = true
		}
		if rule.Keyword != "" && containsKeyword(entry.Message, rule.Keyword) {
			matched = true
		}
		if matched {
			alert := &model.Alert{
				ID:        idgen.Hex(),
				RuleID:    rule.ID,
				SourceID:  entry.SourceID,
				Level:     entry.Level,
				Message:   entry.Message,
				Status:    model.AlertStatusOpen,
				CreatedAt: time.Now(),
			}
			_ = s.store.CreateAlert(alert)
			s.log.Infof("告警触发: rule=%s alert=%s", rule.ID, alert.ID)
		}
	}
}

func levelRank(level string) int {
	switch level {
	case model.LogLevelDebug:
		return 1
	case model.LogLevelInfo:
		return 2
	case model.LogLevelWarn:
		return 3
	case model.LogLevelError:
		return 4
	case model.LogLevelFatal:
		return 5
	}
	return 0
}

func containsKeyword(message, keyword string) bool {
	return strings.Contains(strings.ToLower(message), strings.ToLower(keyword))
}

func (s *Service) GetLogEntry(id string) (*model.LogEntry, error) {
	return s.store.GetLogEntry(id)
}

func (s *Service) ListLogEntries(filter model.LogEntryFilter, page, size int) ([]*model.LogEntry, int, error) {
	all := s.store.ListLogEntries()
	matched := make([]*model.LogEntry, 0, len(all))
	for _, e := range all {
		if filter.Match(e) {
			matched = append(matched, e)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].Timestamp.After(matched[j].Timestamp)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.LogEntry{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) DeleteLogEntriesBySourceID(sourceID string) error {
	if _, err := s.store.GetLogSource(sourceID); err != nil {
		return model.NewValidationError("source_id", "日志源不存在")
	}
	return s.store.DeleteLogEntriesBySourceID(sourceID)
}
