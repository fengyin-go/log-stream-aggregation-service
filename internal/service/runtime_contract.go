package service

import "log-aggregation/internal/store"

type SourceScanResultStreamCoordinator struct { backend *store.SourceScanResultStreamStore }
func NewSourceScanResultStreamCoordinator(b *store.SourceScanResultStreamStore) *SourceScanResultStreamCoordinator { return &SourceScanResultStreamCoordinator{backend: b} }
func (c *SourceScanResultStreamCoordinator) Collect(fail bool) (values []string, err error) { results, errs := c.backend.Stream(fail); for value := range results { values = append(values, value) }; if err := <-errs; err != nil { return values, err }; return values, nil }
