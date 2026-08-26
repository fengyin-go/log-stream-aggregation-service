package store

import "log-aggregation/internal/model"

func (s *MemoryStore) CreateLogSource(src *model.LogSource) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.logSources {
		if exist.Name == src.Name {
			return ErrConflict
		}
	}
	s.logSources[src.ID] = src
	return nil
}

func (s *MemoryStore) GetLogSource(id string) (*model.LogSource, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	src, ok := s.logSources[id]
	if !ok {
		return nil, ErrNotFound
	}
	return src, nil
}

func (s *MemoryStore) GetLogSourceByName(name string) (*model.LogSource, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, src := range s.logSources {
		if src.Name == name {
			return src, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListLogSources() []*model.LogSource {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.LogSource, 0, len(s.logSources))
	for _, src := range s.logSources {
		list = append(list, src)
	}
	return list
}

func (s *MemoryStore) UpdateLogSource(src *model.LogSource) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.logSources[src.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.logSources {
		if exist.ID != src.ID && exist.Name == src.Name {
			return ErrConflict
		}
	}
	s.logSources[src.ID] = src
	return nil
}

func (s *MemoryStore) DeleteLogSource(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.logSources[id]; !ok {
		return ErrNotFound
	}
	delete(s.logSources, id)
	return nil
}
