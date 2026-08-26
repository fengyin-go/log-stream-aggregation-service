package store

import (
    "errors"
    "sync"
)

type RetentionCursorLeasePool struct { mu sync.Mutex; open, limit int }
type RetentionCursorLeaseLease struct { pool *RetentionCursorLeasePool; once sync.Once }
func NewRetentionCursorLeasePool(limit int) *RetentionCursorLeasePool { return &RetentionCursorLeasePool{limit: limit} }
func (p *RetentionCursorLeasePool) Acquire() (*RetentionCursorLeaseLease, error) { p.mu.Lock(); defer p.mu.Unlock(); p.open++; if p.open > p.limit { return nil, errors.New("lease limit") }; return &RetentionCursorLeaseLease{pool: p}, nil }
func (l *RetentionCursorLeaseLease) Close() { l.once.Do(func() { l.pool.mu.Lock(); l.pool.open--; l.pool.mu.Unlock() }) }
func (p *RetentionCursorLeasePool) Open() int { p.mu.Lock(); defer p.mu.Unlock(); return p.open }
