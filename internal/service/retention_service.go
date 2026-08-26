package service

import (
	"sort"
	"time"

	"log-aggregation/internal/model"
	"log-aggregation/pkg/idgen"
)

func (s *Service) CreateRetentionPolicy(input model.RetentionPolicy) (*model.RetentionPolicy, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetLogSource(input.SourceID); err != nil {
		return nil, model.NewValidationError("source_id", "日志源不存在")
	}
	p := &model.RetentionPolicy{
		ID:        idgen.Hex(),
		SourceID:  input.SourceID,
		Policy:    input.Policy,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateRetentionPolicy(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) GetRetentionPolicy(id string) (*model.RetentionPolicy, error) {
	return s.store.GetRetentionPolicy(id)
}

func (s *Service) GetRetentionPolicyBySourceID(sourceID string) (*model.RetentionPolicy, error) {
	return s.store.GetRetentionPolicyBySourceID(sourceID)
}

func (s *Service) ListRetentionPolicies(filter model.RetentionPolicyFilter, page, size int) ([]*model.RetentionPolicy, int, error) {
	all := s.store.ListRetentionPolicies()
	matched := make([]*model.RetentionPolicy, 0, len(all))
	for _, p := range all {
		if filter.Match(p) {
			matched = append(matched, p)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.RetentionPolicy{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateRetentionPolicy(id string, input model.RetentionPolicy) (*model.RetentionPolicy, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	existing, err := s.store.GetRetentionPolicy(id)
	if err != nil {
		return nil, err
	}
	existing.Policy = input.Policy
	if err := s.store.UpdateRetentionPolicy(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteRetentionPolicy(id string) error {
	if _, err := s.store.GetRetentionPolicy(id); err != nil {
		return err
	}
	return s.store.DeleteRetentionPolicy(id)
}

func (s *Service) ApplyRetentionPolicy(sourceID string) (int, error) {
	policy, err := s.store.GetRetentionPolicyBySourceID(sourceID)
	if err != nil {
		return 0, model.NewValidationError("source_id", "未配置保留策略")
	}
	days := model.GetRetentionDays(policy.Policy)
	if days == 0 {
		return 0, nil
	}
	return s.CleanupOldEntries(days * 24)
}
