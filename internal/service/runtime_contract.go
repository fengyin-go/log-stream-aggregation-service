package service

import "log-aggregation/internal/store"

type RetentionCursorLeaseCoordinator struct { pool *store.RetentionCursorLeasePool }
func NewRetentionCursorLeaseCoordinator(p *store.RetentionCursorLeasePool) *RetentionCursorLeaseCoordinator { return &RetentionCursorLeaseCoordinator{pool: p} }
func (c *RetentionCursorLeaseCoordinator) Process(items []string) (processed int, err error) { for range items { lease, err := c.pool.Acquire(); if err != nil { return processed, err }; lease.Close(); processed++ }; return processed, nil }
