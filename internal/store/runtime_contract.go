package store

import (
    "context"
    "sync"
    "time"
)

type SourceHealthProbeStore struct { mu sync.Mutex; first context.Context; delay time.Duration }
func NewSourceHealthProbeStore(delay time.Duration) *SourceHealthProbeStore { return &SourceHealthProbeStore{delay: delay} }
func (s *SourceHealthProbeStore) Wait(ctx context.Context, key string) error {
    s.mu.Lock(); if s.first == nil { s.first = ctx }; active := s.first; s.mu.Unlock()
    timer := time.NewTimer(s.delay)
    defer timer.Stop()
    select { case <-active.Done(): return active.Err(); case <-timer.C: return nil }
}
