package service

import (
	"sort"
	"time"

	"log-aggregation/internal/model"
	"log-aggregation/pkg/idgen"
)

func (s *Service) CreateQuery(input model.Query) (*model.Query, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	q := &model.Query{
		ID:         idgen.Hex(),
		SourceID:   input.SourceID,
		Expression: input.Expression,
		CreatedBy:  input.CreatedBy,
		ExecutedAt: input.ExecutedAt,
		CreatedAt:  time.Now(),
	}
	if err := s.store.CreateQuery(q); err != nil {
		return nil, err
	}
	return q, nil
}

func (s *Service) GetQuery(id string) (*model.Query, error) {
	return s.store.GetQuery(id)
}

func (s *Service) ListQueries(filter model.QueryFilter, page, size int) ([]*model.Query, int, error) {
	all := s.store.ListQueries()
	matched := make([]*model.Query, 0, len(all))
	for _, q := range all {
		if filter.Match(q) {
			matched = append(matched, q)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Query{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateQuery(id string, input model.Query) (*model.Query, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	existing, err := s.store.GetQuery(id)
	if err != nil {
		return nil, err
	}
	existing.SourceID = input.SourceID
	existing.Expression = input.Expression
	existing.CreatedBy = input.CreatedBy
	existing.ExecutedAt = input.ExecutedAt
	if err := s.store.UpdateQuery(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *Service) DeleteQuery(id string) error {
	if _, err := s.store.GetQuery(id); err != nil {
		return err
	}
	return s.store.DeleteQuery(id)
}
