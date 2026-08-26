package store

import (
    "errors"
    "sync"
)

type ExportWriterLeasePool struct { mu sync.Mutex; open, limit int }
type ExportWriterLeaseLease struct { pool *ExportWriterLeasePool; once sync.Once }
func NewExportWriterLeasePool(limit int) *ExportWriterLeasePool { return &ExportWriterLeasePool{limit: limit} }
func (p *ExportWriterLeasePool) Acquire() (*ExportWriterLeaseLease, error) { p.mu.Lock(); defer p.mu.Unlock(); if p.open >= p.limit { return nil, errors.New("lease limit") }; p.open++; return &ExportWriterLeaseLease{pool: p}, nil }
func (l *ExportWriterLeaseLease) Close() { l.once.Do(func() { l.pool.mu.Lock(); l.pool.open--; l.pool.mu.Unlock() }) }
func (p *ExportWriterLeasePool) Open() int { p.mu.Lock(); defer p.mu.Unlock(); return p.open }
