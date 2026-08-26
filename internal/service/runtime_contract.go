package service

import "log-aggregation/internal/store"

type ExportWriterLeaseCoordinator struct { pool *store.ExportWriterLeasePool }
func NewExportWriterLeaseCoordinator(p *store.ExportWriterLeasePool) *ExportWriterLeaseCoordinator { return &ExportWriterLeaseCoordinator{pool: p} }
func (c *ExportWriterLeaseCoordinator) Process(items []string) (processed int, err error) {
	for range items {
		lease, err := c.pool.Acquire()
		if err != nil {
			return processed, err
		}
		processed++
		lease.Close()
	}
	return processed, nil
}
