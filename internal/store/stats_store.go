package store

import (
	"time"
)

func (s *MemoryStore) CountLogSources() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.logSources)
}

func (s *MemoryStore) CountLogEntries() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.logEntries)
}

func (s *MemoryStore) CountLogEntriesBySourceID(sourceID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	count := 0
	for _, e := range s.logEntries {
		if e.SourceID == sourceID {
			count++
		}
	}
	return count
}

func (s *MemoryStore) CountAlerts() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.alerts)
}

func (s *MemoryStore) CountAlertsByStatus(status string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	count := 0
	for _, a := range s.alerts {
		if a.Status == status {
			count++
		}
	}
	return count
}

func (s *MemoryStore) CountLogEntriesByLevel(level string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	count := 0
	for _, e := range s.logEntries {
		if e.Level == level {
			count++
		}
	}
	return count
}

func (s *MemoryStore) CountLogEntriesBySourceAndHour(sourceID string, hour time.Time) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	count := 0
	start := hour.Truncate(time.Hour)
	end := start.Add(time.Hour)
	for _, e := range s.logEntries {
		if e.SourceID == sourceID && !e.Timestamp.Before(start) && e.Timestamp.Before(end) {
			count++
		}
	}
	return count
}

func (s *MemoryStore) CountAlertsBySourceID(sourceID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	count := 0
	for _, a := range s.alerts {
		if a.SourceID == sourceID {
			count++
		}
	}
	return count
}
