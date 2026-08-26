package service

import (
	"sort"
	"time"

	"log-aggregation/internal/model"
	"log-aggregation/pkg/idgen"
)

func (s *Service) CreateAlertRule(input model.AlertRule) (*model.AlertRule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetLogSource(input.SourceID); err != nil {
		return nil, model.NewValidationError("source_id", "日志源不存在")
	}
	rule := &model.AlertRule{
		ID:             idgen.Hex(),
		Name:           input.Name,
		SourceID:       input.SourceID,
		LevelThreshold: input.LevelThreshold,
		Keyword:        input.Keyword,
		Status:         input.Status,
		CreatedAt:      time.Now(),
	}
	if err := s.store.CreateAlertRule(rule); err != nil {
		return nil, err
	}
	return rule, nil
}

func (s *Service) GetAlertRule(id string) (*model.AlertRule, error) {
	return s.store.GetAlertRule(id)
}

func (s *Service) ListAlertRules(filter model.AlertRuleFilter, page, size int) ([]*model.AlertRule, int, error) {
	all := s.store.ListAlertRules()
	matched := make([]*model.AlertRule, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.AlertRule{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateAlertRule(id string, input model.AlertRule) (*model.AlertRule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	existing, err := s.store.GetAlertRule(id)
	if err != nil {
		return nil, err
	}
	existing.Name = input.Name
	existing.LevelThreshold = input.LevelThreshold
	existing.Keyword = input.Keyword
	if err := s.store.UpdateAlertRule(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteAlertRule(id string) error {
	if _, err := s.store.GetAlertRule(id); err != nil {
		return err
	}
	return s.store.DeleteAlertRule(id)
}

func (s *Service) TransitionAlertRuleStatus(id string, toStatus string) (*model.AlertRule, error) {
	existing, err := s.store.GetAlertRule(id)
	if err != nil {
		return nil, err
	}
	if !model.AlertRuleCanTransition(existing.Status, toStatus) {
		return nil, model.NewValidationError("status", "非法的状态流转")
	}
	existing.Status = toStatus
	if err := s.store.UpdateAlertRule(existing); err != nil {
		return nil, err
	}
	return existing, nil
}
