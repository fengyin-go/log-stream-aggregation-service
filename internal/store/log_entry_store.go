package store

import (
	"time"

	"log-aggregation/internal/model"
)

func (s *MemoryStore) CreateLogEntry(e *model.LogEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.logEntries[e.ID] = e
	return nil
}

func (s *MemoryStore) BatchCreateLogEntries(entries []*model.LogEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, e := range entries {
		s.logEntries[e.ID] = e
	}
	return nil
}

func (s *MemoryStore) GetLogEntry(id string) (*model.LogEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.logEntries[id]
	if !ok {
		return nil, ErrNotFound
	}
	return e, nil
}

func (s *MemoryStore) ListLogEntries() []*model.LogEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.LogEntry, 0, len(s.logEntries))
	for _, e := range s.logEntries {
		list = append(list, e)
	}
	return list
}

func (s *MemoryStore) DeleteLogEntry(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.logEntries[id]; !ok {
		return ErrNotFound
	}
	delete(s.logEntries, id)
	return nil
}

func (s *MemoryStore) DeleteLogEntriesBySourceID(sourceID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, e := range s.logEntries {
		if e.SourceID == sourceID {
			delete(s.logEntries, id)
		}
	}
	return nil
}

func (s *MemoryStore) CountLogEntriesBySourceAndLevelAndHour(sourceID, level string, hour time.Time) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	count := 0
	start := hour.Truncate(time.Hour)
	end := start.Add(time.Hour)
	for _, e := range s.logEntries {
		if e.SourceID == sourceID && e.Level == level && !e.Timestamp.Before(start) && e.Timestamp.Before(end) {
			count++
		}
	}
	return count
}
