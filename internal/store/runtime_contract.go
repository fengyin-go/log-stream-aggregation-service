package store

import (
    "context"
    "sync"
    "time"
)

type SourceRecoveryCancellationStore struct { mu sync.Mutex; calls int; delay time.Duration }
func NewSourceRecoveryCancellationStore(delay time.Duration) *SourceRecoveryCancellationStore { return &SourceRecoveryCancellationStore{delay: delay} }
func (s *SourceRecoveryCancellationStore) Attempt(ctx context.Context) error {
    
    s.mu.Lock(); s.calls++; s.mu.Unlock()
    timer := time.NewTimer(s.delay); defer timer.Stop()
    select { case <-ctx.Done(): return ctx.Err(); case <-timer.C: return nil }
}
func (s *SourceRecoveryCancellationStore) Calls() int { s.mu.Lock(); defer s.mu.Unlock(); return s.calls }
