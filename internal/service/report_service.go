package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"log-aggregation/internal/model"
)

type Report struct {
	GeneratedAt      time.Time         `json:"generated_at"`
	SourceCount      int               `json:"source_count"`
	EntryCount       int               `json:"entry_count"`
	AlertCount       int               `json:"alert_count"`
	OpenAlertCount   int               `json:"open_alert_count"`
	LevelDistribution []LevelDistribution `json:"level_distribution"`
	TopSources       []TopSource         `json:"top_sources"`
	RecentErrors     []RecentError       `json:"recent_errors"`
	HealthStatus     string              `json:"health_status"`
}

func (s *Service) GenerateDailyReport() (*Report, error) {
	global := s.GetGlobalStats()
	openAlerts := s.store.CountAlertsByStatus(model.AlertStatusOpen)
	levels := s.GetLevelDistribution("")
	top := s.GetTopSources(5)
	recent := s.GetRecentErrors(10)
	health := "healthy"
	if openAlerts > 50 {
		health = "warning"
	}
	if openAlerts > 100 {
		health = "critical"
	}
	return &Report{
		GeneratedAt:       time.Now(),
		SourceCount:       global.SourceCount,
		EntryCount:        global.EntryCount,
		AlertCount:        global.AlertCount,
		OpenAlertCount:    openAlerts,
		LevelDistribution: levels,
		TopSources:        top,
		RecentErrors:      recent,
		HealthStatus:      health,
	}, nil
}

func (s *Service) GenerateSourceReport(sourceID string) (*Report, error) {
	if _, err := s.store.GetLogSource(sourceID); err != nil {
		return nil, model.NewValidationError("source_id", "日志源不存在")
	}
	entries := s.store.CountLogEntriesBySourceID(sourceID)
	alerts := s.store.CountAlertsBySourceID(sourceID)
	levels := s.GetLevelDistribution(sourceID)
	recent := make([]RecentError, 0)
	all := s.store.ListLogEntries()
	for _, e := range all {
		if e.SourceID == sourceID && (e.Level == model.LogLevelError || e.Level == model.LogLevelFatal) {
			if len(recent) >= 10 {
				break
			}
			recent = append(recent, RecentError{ID: e.ID, SourceID: e.SourceID, Message: e.Message, Timestamp: e.Timestamp})
		}
	}
	health := "healthy"
	if alerts > 10 {
		health = "warning"
	}
	return &Report{
		GeneratedAt:       time.Now(),
		SourceCount:       1,
		EntryCount:        entries,
		AlertCount:        alerts,
		OpenAlertCount:    alerts,
		LevelDistribution: levels,
		TopSources:        []TopSource{{SourceID: sourceID, Count: entries}},
		RecentErrors:      recent,
		HealthStatus:      health,
	}, nil
}

func (s *Service) ExportReportJSON(report *Report) (string, error) {
	b, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (s *Service) ExportReportMarkdown(report *Report) string {
	var sb strings.Builder
	sb.WriteString("# 日志聚合日报\n\n")
	sb.WriteString(fmt.Sprintf("生成时间: %s\n\n", report.GeneratedAt.Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("- 日志源数量: %d\n", report.SourceCount))
	sb.WriteString(fmt.Sprintf("- 日志条目数: %d\n", report.EntryCount))
	sb.WriteString(fmt.Sprintf("- 告警总数: %d\n", report.AlertCount))
	sb.WriteString(fmt.Sprintf("- 未处理告警: %d\n", report.OpenAlertCount))
	sb.WriteString(fmt.Sprintf("- 健康状态: %s\n\n", report.HealthStatus))
	sb.WriteString("## 日志级别分布\n\n")
	sb.WriteString("| 级别 | 数量 | 占比 |\n")
	sb.WriteString("|------|------|------|\n")
	for _, l := range report.LevelDistribution {
		sb.WriteString(fmt.Sprintf("| %s | %d | %.1f%% |\n", l.Level, l.Count, l.Percent))
	}
	sb.WriteString("\n## 活跃日志源 TOP5\n\n")
	sb.WriteString("| 日志源 | 条目数 |\n")
	sb.WriteString("|--------|--------|\n")
	for _, src := range report.TopSources {
		sb.WriteString(fmt.Sprintf("| %s | %d |\n", src.Name, src.Count))
	}
	return sb.String()
}
