package store

import "log-aggregation/internal/model"

func (s *MemoryStore) CreateQuery(q *model.Query) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.queries[q.ID] = q
	return nil
}

func (s *MemoryStore) GetQuery(id string) (*model.Query, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	q, ok := s.queries[id]
	if !ok {
		return nil, ErrNotFound
	}
	return q, nil
}

func (s *MemoryStore) ListQueries() []*model.Query {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Query, 0, len(s.queries))
	for _, q := range s.queries {
		list = append(list, q)
	}
	return list
}

func (s *MemoryStore) UpdateQuery(q *model.Query) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.queries[q.ID]; !ok {
		return ErrNotFound
	}
	s.queries[q.ID] = q
	return nil
}

func (s *MemoryStore) DeleteQuery(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.queries[id]; !ok {
		return ErrNotFound
	}
	delete(s.queries, id)
	return nil
}
