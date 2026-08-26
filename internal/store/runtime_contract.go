package store

import (
    "context"
    "sync"
    "time"
)

type BatchIngestCancellationStore struct { mu sync.Mutex; calls int; delay time.Duration }
func NewBatchIngestCancellationStore(delay time.Duration) *BatchIngestCancellationStore { return &BatchIngestCancellationStore{delay: delay} }
func (s *BatchIngestCancellationStore) Attempt(ctx context.Context) error {
    
    s.mu.Lock(); s.calls++; s.mu.Unlock()
    timer := time.NewTimer(s.delay); defer timer.Stop()
    select { case <-ctx.Done(): return ctx.Err(); case <-timer.C: return nil }
}
func (s *BatchIngestCancellationStore) Calls() int { s.mu.Lock(); defer s.mu.Unlock(); return s.calls }
