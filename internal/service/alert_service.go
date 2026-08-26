package service

import (
	"sort"
	"time"

	"log-aggregation/internal/model"
	"log-aggregation/pkg/idgen"
)

func (s *Service) CreateAlert(input model.Alert) (*model.Alert, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	alert := &model.Alert{
		ID:        idgen.Hex(),
		RuleID:    input.RuleID,
		SourceID:  input.SourceID,
		Level:     input.Level,
		Message:   input.Message,
		Status:    model.AlertStatusOpen,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateAlert(alert); err != nil {
		return nil, err
	}
	return alert, nil
}

func (s *Service) GetAlert(id string) (*model.Alert, error) {
	return s.store.GetAlert(id)
}

func (s *Service) ListAlerts(filter model.AlertFilter, page, size int) ([]*model.Alert, int, error) {
	all := s.store.ListAlerts()
	matched := make([]*model.Alert, 0, len(all))
	for _, a := range all {
		if filter.Match(a) {
			matched = append(matched, a)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Alert{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateAlert(id string, input model.Alert) (*model.Alert, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	existing, err := s.store.GetAlert(id)
	if err != nil {
		return nil, err
	}
	existing.Message = input.Message
	existing.Level = input.Level
	if err := s.store.UpdateAlert(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteAlert(id string) error {
	if _, err := s.store.GetAlert(id); err != nil {
		return err
	}
	return s.store.DeleteAlert(id)
}

func (s *Service) TransitionAlertStatus(id string, toStatus string) (*model.Alert, error) {
	existing, err := s.store.GetAlert(id)
	if err != nil {
		return nil, err
	}
	if !model.AlertCanTransition(existing.Status, toStatus) {
		return nil, model.NewValidationError("status", "非法的状态流转")
	}
	existing.Status = toStatus
	if err := s.store.UpdateAlert(existing); err != nil {
		return nil, err
	}
	return existing, nil
}
