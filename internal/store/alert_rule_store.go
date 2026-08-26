package store

import "log-aggregation/internal/model"

func (s *MemoryStore) CreateAlertRule(r *model.AlertRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.alertRules[r.ID] = r
	return nil
}

func (s *MemoryStore) GetAlertRule(id string) (*model.AlertRule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.alertRules[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) ListAlertRules() []*model.AlertRule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.AlertRule, 0, len(s.alertRules))
	for _, r := range s.alertRules {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) UpdateAlertRule(r *model.AlertRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.alertRules[r.ID]; !ok {
		return ErrNotFound
	}
	s.alertRules[r.ID] = r
	return nil
}

func (s *MemoryStore) DeleteAlertRule(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.alertRules[id]; !ok {
		return ErrNotFound
	}
	delete(s.alertRules, id)
	return nil
}

func (s *MemoryStore) ListActiveAlertRulesBySourceID(sourceID string) []*model.AlertRule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.AlertRule, 0)
	for _, r := range s.alertRules {
		if r.SourceID == sourceID && r.Status == model.AlertRuleStatusActive {
			list = append(list, r)
		}
	}
	return list
}
