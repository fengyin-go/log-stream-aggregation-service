package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"log-aggregation/internal/model"
)

type CSVExportResult struct {
	ContentType string `json:"-"`
	Filename    string `json:"-"`
	Body        string `json:"-"`
}

func (s *Service) ExportLogEntriesCSV(sourceID string) (*CSVExportResult, error) {
	if _, err := s.store.GetLogSource(sourceID); err != nil {
		return nil, model.NewValidationError("source_id", "日志源不存在")
	}
	entries, _, err := s.ListLogEntries(model.LogEntryFilter{SourceID: sourceID}, 1, 1000)
	if err != nil {
		return nil, err
	}
	var sb strings.Builder
	sb.WriteString("id,source_id,level,message,timestamp,tags\n")
	for _, e := range entries {
		tags := strings.Join(e.Tags, "|")
		line := fmt.Sprintf("%s,%s,%s,%q,%s,%s\n", e.ID, e.SourceID, e.Level, e.Message, e.Timestamp.Format("2006-01-02T15:04:05"), tags)
		sb.WriteString(line)
	}
	return &CSVExportResult{
		ContentType: "text/csv",
		Filename:    fmt.Sprintf("logs_%s.csv", sourceID),
		Body:        sb.String(),
	}, nil
}

func (s *Service) ExportAlertsCSV(sourceID string) (*CSVExportResult, error) {
	all := s.store.ListAlerts()
	var sb strings.Builder
	sb.WriteString("id,rule_id,source_id,level,message,status,created_at\n")
	for _, a := range all {
		if sourceID != "" && a.SourceID != sourceID {
			continue
		}
		line := fmt.Sprintf("%s,%s,%s,%s,%q,%s,%s\n", a.ID, a.RuleID, a.SourceID, a.Level, a.Message, a.Status, a.CreatedAt.Format("2006-01-02T15:04:05"))
		sb.WriteString(line)
	}
	return &CSVExportResult{
		ContentType: "text/csv",
		Filename:    "alerts.csv",
		Body:        sb.String(),
	}, nil
}

func (s *Service) ExportSourcesJSON() (string, error) {
	sources := s.store.ListLogSources()
	result := make([]map[string]interface{}, 0, len(sources))
	for _, src := range sources {
		result = append(result, map[string]interface{}{
			"id":         src.ID,
			"name":       src.Name,
			"type":       src.Type,
			"host":       src.Host,
			"path":       src.Path,
			"status":     src.Status,
			"entry_count": s.store.CountLogEntriesBySourceID(src.ID),
			"alert_count": s.store.CountAlertsBySourceID(src.ID),
			"created_at": src.CreatedAt,
		})
	}
	b, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
