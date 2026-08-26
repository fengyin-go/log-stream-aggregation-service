package service

import (
	"sort"
	"time"

	"log-aggregation/internal/model"
	"log-aggregation/pkg/idgen"
)

func (s *Service) CreateLogSource(input model.LogSource) (*model.LogSource, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	src := &model.LogSource{
		ID:        idgen.Hex(),
		Name:      input.Name,
		Type:      input.Type,
		Host:      input.Host,
		Path:      input.Path,
		Status:    input.Status,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateLogSource(src); err != nil {
		return nil, err
	}
	s.log.Infof("创建日志源: %s", src.Name)
	return src, nil
}

func (s *Service) GetLogSource(id string) (*model.LogSource, error) {
	return s.store.GetLogSource(id)
}

func (s *Service) ListLogSources(filter model.LogSourceFilter, page, size int) ([]*model.LogSource, int, error) {
	all := s.store.ListLogSources()
	matched := make([]*model.LogSource, 0, len(all))
	for _, src := range all {
		if filter.Match(src) {
			matched = append(matched, src)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.LogSource{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateLogSource(id string, input model.LogSource) (*model.LogSource, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	existing, err := s.store.GetLogSource(id)
	if err != nil {
		return nil, err
	}
	existing.Name = input.Name
	existing.Type = input.Type
	existing.Host = input.Host
	existing.Path = input.Path
	if err := s.store.UpdateLogSource(existing); err != nil {
		return nil, err
	}
	s.log.Infof("更新日志源: %s", existing.Name)
	return existing, nil
}

func (s *Service) DeleteLogSource(id string) error {
	if _, err := s.store.GetLogSource(id); err != nil {
		return err
	}
	if err := s.store.DeleteLogEntriesBySourceID(id); err != nil {
		return err
	}
	if err := s.store.DeleteLogSource(id); err != nil {
		return err
	}
	s.log.Infof("删除日志源: %s", id)
	return nil
}

func (s *Service) TransitionLogSourceStatus(id string, toStatus string) (*model.LogSource, error) {
	existing, err := s.store.GetLogSource(id)
	if err != nil {
		return nil, err
	}
	if !model.LogSourceCanTransition(existing.Status, toStatus) {
		return nil, model.NewValidationError("status", "非法的状态流转")
	}
	existing.Status = toStatus
	if err := s.store.UpdateLogSource(existing); err != nil {
		return nil, err
	}
	s.log.Infof("日志源状态变更: %s -> %s", id, toStatus)
	return existing, nil
}
