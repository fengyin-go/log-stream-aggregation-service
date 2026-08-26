package store

import "log-aggregation/internal/model"

func (s *MemoryStore) CreateRetentionPolicy(p *model.RetentionPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.retentionPolicies {
		if exist.SourceID == p.SourceID {
			return ErrConflict
		}
	}
	s.retentionPolicies[p.ID] = p
	return nil
}

func (s *MemoryStore) GetRetentionPolicy(id string) (*model.RetentionPolicy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.retentionPolicies[id]
	if !ok {
		return nil, ErrNotFound
	}
	return p, nil
}

func (s *MemoryStore) GetRetentionPolicyBySourceID(sourceID string) (*model.RetentionPolicy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.retentionPolicies {
		if p.SourceID == sourceID {
			return p, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListRetentionPolicies() []*model.RetentionPolicy {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.RetentionPolicy, 0, len(s.retentionPolicies))
	for _, p := range s.retentionPolicies {
		list = append(list, p)
	}
	return list
}

func (s *MemoryStore) UpdateRetentionPolicy(p *model.RetentionPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.retentionPolicies[p.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.retentionPolicies {
		if exist.ID != p.ID && exist.SourceID == p.SourceID {
			return ErrConflict
		}
	}
	s.retentionPolicies[p.ID] = p
	return nil
}

func (s *MemoryStore) DeleteRetentionPolicy(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.retentionPolicies[id]; !ok {
		return ErrNotFound
	}
	delete(s.retentionPolicies, id)
	return nil
}
