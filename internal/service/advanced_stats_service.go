package service

import (
	"sort"
	"time"

	"log-aggregation/internal/model"
)

type TopSource struct {
	SourceID string `json:"source_id"`
	Name     string `json:"name"`
	Count    int    `json:"count"`
}

type RecentError struct {
	ID        string    `json:"id"`
	SourceID  string    `json:"source_id"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type LevelDistribution struct {
	Level   string  `json:"level"`
	Count   int     `json:"count"`
	Percent float64 `json:"percent"`
}

type HourlyPeak struct {
	Hour  time.Time `json:"hour"`
	Count int       `json:"count"`
}

func (s *Service) GetTopSources(limit int) []TopSource {
	if limit <= 0 {
		limit = 5
	}
	sources := s.store.ListLogSources()
	result := make([]TopSource, 0, len(sources))
	for _, src := range sources {
		result = append(result, TopSource{SourceID: src.ID, Name: src.Name, Count: s.store.CountLogEntriesBySourceID(src.ID)})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})
	if len(result) > limit {
		result = result[:limit]
	}
	return result
}

func (s *Service) GetRecentErrors(limit int) []RecentError {
	if limit <= 0 {
		limit = 10
	}
	all := s.store.ListLogEntries()
	errors := make([]*model.LogEntry, 0)
	for _, e := range all {
		if e.Level == model.LogLevelError || e.Level == model.LogLevelFatal {
			errors = append(errors, e)
		}
	}
	sort.Slice(errors, func(i, j int) bool {
		return errors[i].Timestamp.After(errors[j].Timestamp)
	})
	result := make([]RecentError, 0, limit)
	for i := 0; i < len(errors) && i < limit; i++ {
		result = append(result, RecentError{ID: errors[i].ID, SourceID: errors[i].SourceID, Message: errors[i].Message, Timestamp: errors[i].Timestamp})
	}
	return result
}

func (s *Service) GetLevelDistribution(sourceID string) []LevelDistribution {
	all := s.store.ListLogEntries()
	counts := map[string]int{}
	total := 0
	for _, e := range all {
		if sourceID != "" && e.SourceID != sourceID {
			continue
		}
		counts[e.Level]++
		total++
	}
	levels := []string{model.LogLevelDebug, model.LogLevelInfo, model.LogLevelWarn, model.LogLevelError, model.LogLevelFatal}
	result := make([]LevelDistribution, 0, len(levels))
	for _, level := range levels {
		percent := 0.0
		if total > 0 {
			percent = float64(counts[level]) / float64(total) * 100
		}
		result = append(result, LevelDistribution{Level: level, Count: counts[level], Percent: percent})
	}
	return result
}

func (s *Service) GetHourlyPeak(sourceID string, hours int) []HourlyPeak {
	if hours <= 0 {
		hours = 24
	}
	if hours > 168 {
		hours = 168
	}
	now := time.Now()
	result := make([]HourlyPeak, 0, hours)
	for i := hours - 1; i >= 0; i-- {
		hour := now.Add(-time.Duration(i) * time.Hour).Truncate(time.Hour)
		count := s.store.CountLogEntriesBySourceAndHour(sourceID, hour)
		result = append(result, HourlyPeak{Hour: hour, Count: count})
	}
	return result
}
