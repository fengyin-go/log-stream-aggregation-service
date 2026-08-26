package service

import "log-aggregation/internal/store"

type IndexReaderLeaseCoordinator struct { pool *store.IndexReaderLeasePool }
func NewIndexReaderLeaseCoordinator(p *store.IndexReaderLeasePool) *IndexReaderLeaseCoordinator { return &IndexReaderLeaseCoordinator{pool: p} }
func (c *IndexReaderLeaseCoordinator) Process(items []string) (processed int, err error) {
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
