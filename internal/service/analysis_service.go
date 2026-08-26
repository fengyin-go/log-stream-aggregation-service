package service

import (
	"sort"
	"strings"
	"time"

	"log-aggregation/internal/model"
)

func (s *Service) AnalyzeLogSource(sourceID string, hours int) (*model.LogAnalysis, error) {
	if _, err := s.store.GetLogSource(sourceID); err != nil {
		return nil, model.NewValidationError("source_id", "日志源不存在")
	}
	if hours <= 0 {
		hours = 24
	}
	if hours > 168 {
		hours = 168
	}
	end := time.Now()
	start := end.Add(-time.Duration(hours) * time.Hour)

	all := s.store.ListLogEntries()
	analysis := &model.LogAnalysis{
		SourceID:    sourceID,
		PeriodStart: start,
		PeriodEnd:   end,
	}
	keywordCount := make(map[string]int)
	for _, e := range all {
		if e.SourceID != sourceID {
			continue
		}
		if e.Timestamp.Before(start) || e.Timestamp.After(end) {
			continue
		}
		analysis.TotalCount++
		switch e.Level {
		case model.LogLevelDebug:
			analysis.DebugCount++
		case model.LogLevelInfo:
			analysis.InfoCount++
		case model.LogLevelWarn:
			analysis.WarnCount++
		case model.LogLevelError:
			analysis.ErrorCount++
		case model.LogLevelFatal:
			analysis.FatalCount++
		}
		words := strings.Fields(e.Message)
		for _, w := range words {
			w = strings.ToLower(strings.TrimSpace(w))
			if len(w) > 3 && !isStopWord(w) {
				keywordCount[w]++
			}
		}
	}
	if analysis.TotalCount > 0 {
		analysis.ErrorRate = float64(analysis.ErrorCount+analysis.FatalCount) / float64(analysis.TotalCount)
	}
	var kfs []model.KeywordFreq
	for k, c := range keywordCount {
		kfs = append(kfs, model.KeywordFreq{Keyword: k, Count: c})
	}
	sort.Slice(kfs, func(i, j int) bool {
		return kfs[i].Count > kfs[j].Count
	})
	if len(kfs) > 10 {
		kfs = kfs[:10]
	}
	analysis.TopKeywords = kfs
	return analysis, nil
}

func isStopWord(w string) bool {
	stopWords := map[string]bool{
		"the": true, "this": true, "that": true, "with": true, "from": true,
		"have": true, "been": true, "were": true, "they": true, "will": true,
	}
	return stopWords[w]
}

func (s *Service) GetSourceHealth(sourceID string) (*model.SourceHealth, error) {
	src, err := s.store.GetLogSource(sourceID)
	if err != nil {
		return nil, model.NewValidationError("source_id", "日志源不存在")
	}
	now := time.Now()
	start := now.Add(-24 * time.Hour)
	all := s.store.ListLogEntries()
	entryCount24h := 0
	var lastEntryTime time.Time
	for _, e := range all {
		if e.SourceID != sourceID {
			continue
		}
		if e.Timestamp.After(lastEntryTime) {
			lastEntryTime = e.Timestamp
		}
		if e.Timestamp.After(start) && e.Timestamp.Before(now) {
			entryCount24h++
		}
	}
	alertCount24h := 0
	allAlerts := s.store.ListAlerts()
	for _, a := range allAlerts {
		if a.SourceID == sourceID && a.CreatedAt.After(start) && a.CreatedAt.Before(now) {
			alertCount24h++
		}
	}
	healthy := src.Status == model.LogSourceStatusActive && entryCount24h > 0
	return &model.SourceHealth{
		SourceID:      sourceID,
		SourceName:    src.Name,
		Status:        src.Status,
		LastEntryTime: lastEntryTime,
		EntryCount24h: entryCount24h,
		AlertCount24h: alertCount24h,
		Healthy:       healthy,
	}, nil
}

func (s *Service) ListSourceHealth() []*model.SourceHealth {
	sources := s.store.ListLogSources()
	result := make([]*model.SourceHealth, 0, len(sources))
	for _, src := range sources {
		h, _ := s.GetSourceHealth(src.ID)
		if h != nil {
			result = append(result, h)
		}
	}
	return result
}

func (s *Service) GetErrorRateTrend(sourceID string, hours int) []model.TimeTrend {
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
		total := s.store.CountLogEntriesBySourceAndHour(sourceID, hour)
		errors := s.store.CountLogEntriesBySourceAndLevelAndHour(sourceID, model.LogLevelError, hour)
		result = append(result, model.TimeTrend{Time: hour, Count: errors})
		if total > 0 {
			result[len(result)-1].Count = int(float64(errors) / float64(total) * 100)
		}
	}
	return result
}
