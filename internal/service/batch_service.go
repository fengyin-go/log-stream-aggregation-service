package service

import (
	"log-aggregation/internal/model"
	"log-aggregation/pkg/idgen"
	"time"
)

func (s *Service) BatchAcknowledgeAlerts(alertIDs []string) (int, error) {
	if len(alertIDs) == 0 {
		return 0, model.NewValidationError("alert_ids", "告警 ID 列表不能为空")
	}
	count := 0
	for _, id := range alertIDs {
		alert, err := s.store.GetAlert(id)
		if err != nil {
			continue
		}
		if alert.Status != model.AlertStatusOpen {
			continue
		}
		alert.Status = model.AlertStatusAcknowledged
		if err := s.store.UpdateAlert(alert); err == nil {
			count++
		}
	}
	s.log.Infof("批量确认告警: %d/%d", count, len(alertIDs))
	return count, nil
}

func (s *Service) BatchResolveAlerts(alertIDs []string) (int, error) {
	if len(alertIDs) == 0 {
		return 0, model.NewValidationError("alert_ids", "告警 ID 列表不能为空")
	}
	count := 0
	for _, id := range alertIDs {
		alert, err := s.store.GetAlert(id)
		if err != nil {
			continue
		}
		if alert.Status == model.AlertStatusResolved {
			continue
		}
		alert.Status = model.AlertStatusResolved
		if err := s.store.UpdateAlert(alert); err == nil {
			count++
		}
	}
	s.log.Infof("批量解决告警: %d/%d", count, len(alertIDs))
	return count, nil
}

func (s *Service) BatchDeleteLogEntries(entryIDs []string) (int, error) {
	if len(entryIDs) == 0 {
		return 0, model.NewValidationError("entry_ids", "日志 ID 列表不能为空")
	}
	count := 0
	for _, id := range entryIDs {
		if err := s.store.DeleteLogEntry(id); err == nil {
			count++
		}
	}
	s.log.Infof("批量删除日志: %d/%d", count, len(entryIDs))
	return count, nil
}

func (s *Service) BatchPauseSources(sourceIDs []string) (int, error) {
	if len(sourceIDs) == 0 {
		return 0, model.NewValidationError("source_ids", "日志源 ID 列表不能为空")
	}
	count := 0
	for _, id := range sourceIDs {
		src, err := s.store.GetLogSource(id)
		if err != nil {
			continue
		}
		if src.Status == model.LogSourceStatusPaused {
			continue
		}
		src.Status = model.LogSourceStatusPaused
		if err := s.store.UpdateLogSource(src); err == nil {
			count++
		}
	}
	s.log.Infof("批量暂停日志源: %d/%d", count, len(sourceIDs))
	return count, nil
}

func (s *Service) BatchResumeSources(sourceIDs []string) (int, error) {
	if len(sourceIDs) == 0 {
		return 0, model.NewValidationError("source_ids", "日志源 ID 列表不能为空")
	}
	count := 0
	for _, id := range sourceIDs {
		src, err := s.store.GetLogSource(id)
		if err != nil {
			continue
		}
		if src.Status == model.LogSourceStatusActive {
			continue
		}
		src.Status = model.LogSourceStatusActive
		if err := s.store.UpdateLogSource(src); err == nil {
			count++
		}
	}
	s.log.Infof("批量恢复日志源: %d/%d", count, len(sourceIDs))
	return count, nil
}

func (s *Service) CleanupOldEntries(hours int) (int, error) {
	if hours <= 0 {
		return 0, model.NewValidationError("hours", "清理时长必须大于 0")
	}
	cutoff := time.Now().Add(-time.Duration(hours) * time.Hour)
	all := s.store.ListLogEntries()
	count := 0
	for _, e := range all {
		if e.Timestamp.Before(cutoff) {
			if err := s.store.DeleteLogEntry(e.ID); err == nil {
				count++
			}
		}
	}
	s.log.Infof("清理旧日志: %d 条", count)
	return count, nil
}

func (s *Service) DuplicateAlertRule(ruleID string) (*model.AlertRule, error) {
	rule, err := s.store.GetAlertRule(ruleID)
	if err != nil {
		return nil, err
	}
	newRule := &model.AlertRule{
		ID:             idgen.Hex(),
		Name:           rule.Name + " (copy)",
		SourceID:       rule.SourceID,
		LevelThreshold: rule.LevelThreshold,
		Keyword:        rule.Keyword,
		Status:         model.AlertRuleStatusActive,
		CreatedAt:      time.Now(),
	}
	if err := s.store.CreateAlertRule(newRule); err != nil {
		return nil, err
	}
	return newRule, nil
}
