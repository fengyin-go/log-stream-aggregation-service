package store

import (
    "context"
    "sync"
    "time"
)

type RetentionSweepCancellationStore struct { mu sync.Mutex; calls int; delay time.Duration }
func NewRetentionSweepCancellationStore(delay time.Duration) *RetentionSweepCancellationStore { return &RetentionSweepCancellationStore{delay: delay} }
func (s *RetentionSweepCancellationStore) Attempt(ctx context.Context) error {
    
    s.mu.Lock(); s.calls++; s.mu.Unlock()
    timer := time.NewTimer(s.delay); defer timer.Stop()
    select { case <-ctx.Done(): return ctx.Err(); case <-timer.C: return nil }
}
func (s *RetentionSweepCancellationStore) Calls() int { s.mu.Lock(); defer s.mu.Unlock(); return s.calls }
