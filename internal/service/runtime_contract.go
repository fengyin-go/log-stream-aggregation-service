package service

import "log-aggregation/internal/store"

type HealthCheckResultStreamCoordinator struct { backend *store.HealthCheckResultStreamStore }
func NewHealthCheckResultStreamCoordinator(b *store.HealthCheckResultStreamStore) *HealthCheckResultStreamCoordinator { return &HealthCheckResultStreamCoordinator{backend: b} }
func (c *HealthCheckResultStreamCoordinator) Collect(fail bool) (values []string, err error) { results, errs := c.backend.Stream(fail); for value := range results { values = append(values, value) }; if err := <-errs; err != nil { return values, err }; return values, nil }
