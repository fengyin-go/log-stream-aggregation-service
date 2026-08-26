package model

import "time"

type TimeTrend struct {
	Time  time.Time `json:"time"`
	Count int       `json:"count"`
}

type LogAnalysis struct {
	SourceID    string    `json:"source_id"`
	TotalCount  int       `json:"total_count"`
	ErrorCount  int       `json:"error_count"`
	WarnCount   int       `json:"warn_count"`
	InfoCount   int       `json:"info_count"`
	FatalCount  int       `json:"fatal_count"`
	DebugCount  int       `json:"debug_count"`
	ErrorRate   float64   `json:"error_rate"`
	TopKeywords []KeywordFreq `json:"top_keywords"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
}

type KeywordFreq struct {
	Keyword string `json:"keyword"`
	Count   int    `json:"count"`
}

type SourceHealth struct {
	SourceID       string    `json:"source_id"`
	SourceName     string    `json:"source_name"`
	Status         string    `json:"status"`
	LastEntryTime  time.Time `json:"last_entry_time"`
	EntryCount24h  int       `json:"entry_count_24h"`
	AlertCount24h  int       `json:"alert_count_24h"`
	Healthy        bool      `json:"healthy"`
}
