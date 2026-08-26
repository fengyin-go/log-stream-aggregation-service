package service

import (
	"sort"
	"time"

	"log-aggregation/internal/model"
	"log-aggregation/pkg/idgen"
)

func (s *Service) CreateIndex(input model.Index) (*model.Index, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetLogSource(input.SourceID); err != nil {
		return nil, model.NewValidationError("source_id", "日志源不存在")
	}
	idx := &model.Index{
		ID:        idgen.Hex(),
		Name:      input.Name,
		SourceID:  input.SourceID,
		Field:     input.Field,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateIndex(idx); err != nil {
		return nil, err
	}
	return idx, nil
}

func (s *Service) GetIndex(id string) (*model.Index, error) {
	return s.store.GetIndex(id)
}

func (s *Service) ListIndexes(filter model.IndexFilter, page, size int) ([]*model.Index, int, error) {
	all := s.store.ListIndexes()
	matched := make([]*model.Index, 0, len(all))
	for _, i := range all {
		if filter.Match(i) {
			matched = append(matched, i)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Index{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateIndex(id string, input model.Index) (*model.Index, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	existing, err := s.store.GetIndex(id)
	if err != nil {
		return nil, err
	}
	existing.Name = input.Name
	existing.Field = input.Field
	if err := s.store.UpdateIndex(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteIndex(id string) error {
	if _, err := s.store.GetIndex(id); err != nil {
		return err
	}
	return s.store.DeleteIndex(id)
}
