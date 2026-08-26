package service

import (
	"time"

	"log-aggregation/internal/model"
)

type HealthCheckResult struct {
	SourceID      string    `json:"source_id"`
	SourceName    string    `json:"source_name"`
	Reachable     bool      `json:"reachable"`
	ResponseTime  int64     `json:"response_time_ms"`
	LastCheckedAt time.Time `json:"last_checked_at"`
	Message       string    `json:"message"`
}

func (s *Service) CheckSourceHealth(sourceID string) (*HealthCheckResult, error) {
	src, err := s.store.GetLogSource(sourceID)
	if err != nil {
		return nil, err
	}
	result := &HealthCheckResult{
		SourceID:      src.ID,
		SourceName:    src.Name,
		LastCheckedAt: time.Now(),
	}
	if src.Status == model.LogSourceStatusPaused {
		result.Reachable = false
		result.Message = "日志源已暂停"
		return result, nil
	}
	start := time.Now()
	count := s.store.CountLogEntriesBySourceID(sourceID)
	elapsed := time.Since(start)
	result.ResponseTime = elapsed.Milliseconds()
	result.Reachable = true
	if count == 0 {
		result.Message = "可达，但最近 24 小时无日志"
	} else {
		result.Message = "健康"
	}
	return result, nil
}

func (s *Service) CheckAllSourcesHealth() []*HealthCheckResult {
	sources := s.store.ListLogSources()
	results := make([]*HealthCheckResult, 0, len(sources))
	for _, src := range sources {
		r, _ := s.CheckSourceHealth(src.ID)
		if r != nil {
			results = append(results, r)
		}
	}
	return results
}

type SystemHealth struct {
	SourcesTotal   int `json:"sources_total"`
	SourcesActive  int `json:"sources_active"`
	SourcesPaused  int `json:"sources_paused"`
	EntriesTotal   int `json:"entries_total"`
	AlertsOpen     int `json:"alerts_open"`
	AlertsToday    int `json:"alerts_today"`
	Healthy        bool `json:"healthy"`
}

func (s *Service) GetSystemHealth() SystemHealth {
	sources := s.store.ListLogSources()
	active := 0
	paused := 0
	for _, src := range sources {
		if src.Status == model.LogSourceStatusActive {
			active++
		} else {
			paused++
		}
	}
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	alertsToday := 0
	for _, a := range s.store.ListAlerts() {
		if a.CreatedAt.After(startOfDay) {
			alertsToday++
		}
	}
	healthy := active > 0 && s.store.CountAlertsByStatus(model.AlertStatusOpen) < 100
	return SystemHealth{
		SourcesTotal:  len(sources),
		SourcesActive: active,
		SourcesPaused: paused,
		EntriesTotal:  s.store.CountLogEntries(),
		AlertsOpen:    s.store.CountAlertsByStatus(model.AlertStatusOpen),
		AlertsToday:   alertsToday,
		Healthy:       healthy,
	}
}
