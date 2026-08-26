package service

import "log-aggregation/internal/store"

type ExportChunkResultStreamCoordinator struct { backend *store.ExportChunkResultStreamStore }
func NewExportChunkResultStreamCoordinator(b *store.ExportChunkResultStreamStore) *ExportChunkResultStreamCoordinator { return &ExportChunkResultStreamCoordinator{backend: b} }
func (c *ExportChunkResultStreamCoordinator) Collect(fail bool) (values []string, err error) { results, errs := c.backend.Stream(fail); for value := range results { values = append(values, value) }; if err := <-errs; err != nil { return values, err }; return values, nil }
