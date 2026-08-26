package store

import (
    "context"
    "sync"
    "time"
)

type TailCursorProbeStore struct { mu sync.Mutex; first context.Context; delay time.Duration }
func NewTailCursorProbeStore(delay time.Duration) *TailCursorProbeStore { return &TailCursorProbeStore{delay: delay} }
func (s *TailCursorProbeStore) Wait(ctx context.Context, key string) error {
    s.mu.Lock(); if s.first == nil { s.first = ctx }; active := s.first; s.mu.Unlock()
    timer := time.NewTimer(s.delay)
    defer timer.Stop()
    select { case <-active.Done(): return active.Err(); case <-timer.C: return nil }
}
