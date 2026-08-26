package store

import "log-aggregation/internal/model"

func (s *MemoryStore) CreateIndex(i *model.Index) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.indexes[i.ID] = i
	return nil
}

func (s *MemoryStore) GetIndex(id string) (*model.Index, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	i, ok := s.indexes[id]
	if !ok {
		return nil, ErrNotFound
	}
	return i, nil
}

func (s *MemoryStore) ListIndexes() []*model.Index {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Index, 0, len(s.indexes))
	for _, i := range s.indexes {
		list = append(list, i)
	}
	return list
}

func (s *MemoryStore) UpdateIndex(i *model.Index) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.indexes[i.ID]; !ok {
		return ErrNotFound
	}
	s.indexes[i.ID] = i
	return nil
}

func (s *MemoryStore) DeleteIndex(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.indexes[id]; !ok {
		return ErrNotFound
	}
	delete(s.indexes, id)
	return nil
}
