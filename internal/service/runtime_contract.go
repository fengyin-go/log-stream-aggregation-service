package service

import (
    "errors"
    "log-aggregation/internal/store"
)

type LogMetadataAssemblyCoordinator struct { backend *store.LogMetadataAssemblyStore }
func NewLogMetadataAssemblyCoordinator(b *store.LogMetadataAssemblyStore) *LogMetadataAssemblyCoordinator { return &LogMetadataAssemblyCoordinator{backend: b} }
func (c *LogMetadataAssemblyCoordinator) Build(key string, fail bool) (item *store.Assembly, err error) { defer func() { if recover() != nil { item, _ = c.backend.Get(key); err = errors.New("assembly failed") } }(); item = c.backend.Build(key, fail); return item, nil }
