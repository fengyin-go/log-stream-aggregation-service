package service

import (
	"time"

	"log-aggregation/internal/model"
)

type GlobalStats struct {
	SourceCount int `json:"source_count"`
	EntryCount  int `json:"entry_count"`
	AlertCount  int `json:"alert_count"`
}

type LevelStats struct {
	Level string `json:"level"`
	Count int    `json:"count"`
}

type SourceStats struct {
	SourceID string `json:"source_id"`
	Count    int    `json:"count"`
}

type AlertStatusStats struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

func (s *Service) GetGlobalStats() GlobalStats {
	return GlobalStats{
		SourceCount: s.store.CountLogSources(),
		EntryCount:  s.store.CountLogEntries(),
		AlertCount:  s.store.CountAlerts(),
	}
}

func (s *Service) GetLevelStats() []LevelStats {
	levels := []string{model.LogLevelDebug, model.LogLevelInfo, model.LogLevelWarn, model.LogLevelError, model.LogLevelFatal}
	result := make([]LevelStats, 0, len(levels))
	for _, level := range levels {
		result = append(result, LevelStats{Level: level, Count: s.store.CountLogEntriesByLevel(level)})
	}
	return result
}

func (s *Service) GetSourceStats() []SourceStats {
	sources := s.store.ListLogSources()
	result := make([]SourceStats, 0, len(sources))
	for _, src := range sources {
		result = append(result, SourceStats{SourceID: src.ID, Count: s.store.CountLogEntriesBySourceID(src.ID)})
	}
	return result
}

func (s *Service) GetTimeTrend(sourceID string, hours int) []model.TimeTrend {
	if hours <= 0 {
		hours = 24
	}
	if hours > 168 {
		hours = 168
	}
	now := time.Now()
	result := make([]model.TimeTrend, 0, hours)
	for i := hours - 1; i >= 0; i-- {
		hour := now.Add(-time.Duration(i) * time.Hour).Truncate(time.Hour)
		count := s.store.CountLogEntriesBySourceAndHour(sourceID, hour)
		result = append(result, model.TimeTrend{Time: hour, Count: count})
	}
	return result
}

func (s *Service) GetAlertStatusStats() []AlertStatusStats {
	statuses := []string{model.AlertStatusOpen, model.AlertStatusAcknowledged, model.AlertStatusResolved}
	result := make([]AlertStatusStats, 0, len(statuses))
	for _, status := range statuses {
		result = append(result, AlertStatusStats{Status: status, Count: s.store.CountAlertsByStatus(status)})
	}
	return result
}

func (s *Service) ExportLogEntriesSummary(sourceID string) ([]map[string]interface{}, error) {
	if _, err := s.store.GetLogSource(sourceID); err != nil {
		return nil, model.NewValidationError("source_id", "日志源不存在")
	}
	entries, _, err := s.ListLogEntries(model.LogEntryFilter{SourceID: sourceID}, 1, 1000)
	if err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, 0, len(entries))
	for _, e := range entries {
		result = append(result, map[string]interface{}{
			"id":        e.ID,
			"level":     e.Level,
			"message":   e.Message,
			"timestamp": e.Timestamp,
			"tags":      e.Tags,
		})
	}
	return result, nil
}
