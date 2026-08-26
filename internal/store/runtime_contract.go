package store

import (
    "errors"
    "sync"
)

type IndexReaderLeasePool struct { mu sync.Mutex; open, limit int }
type IndexReaderLeaseLease struct { pool *IndexReaderLeasePool; once sync.Once }
func NewIndexReaderLeasePool(limit int) *IndexReaderLeasePool { return &IndexReaderLeasePool{limit: limit} }
func (p *IndexReaderLeasePool) Acquire() (*IndexReaderLeaseLease, error) { p.mu.Lock(); defer p.mu.Unlock(); if p.open >= p.limit { return nil, errors.New("lease limit") }; p.open++; return &IndexReaderLeaseLease{pool: p}, nil }
func (l *IndexReaderLeaseLease) Close() { l.once.Do(func() { l.pool.mu.Lock(); l.pool.open--; l.pool.mu.Unlock() }) }
func (p *IndexReaderLeasePool) Open() int { p.mu.Lock(); defer p.mu.Unlock(); return p.open }
