package service

import (
    "errors"
    "log-aggregation/internal/store"
)

type ParsedEntryAssemblyCoordinator struct { backend *store.ParsedEntryAssemblyStore }
func NewParsedEntryAssemblyCoordinator(b *store.ParsedEntryAssemblyStore) *ParsedEntryAssemblyCoordinator { return &ParsedEntryAssemblyCoordinator{backend: b} }
func (c *ParsedEntryAssemblyCoordinator) Build(key string, fail bool) (item *store.Assembly, err error) { defer func() { if recover() != nil { item, _ = c.backend.Get(key); err = errors.New("assembly failed") } }(); item = c.backend.Build(key, fail); return item, nil }
