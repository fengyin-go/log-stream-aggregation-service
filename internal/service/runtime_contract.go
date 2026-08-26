package service

import "log-aggregation/internal/store"

type QueryTokenArchiveCoordinator struct { backend *store.QueryTokenArchiveStore; pending []byte }
func NewQueryTokenArchiveCoordinator(b *store.QueryTokenArchiveStore) *QueryTokenArchiveCoordinator { return &QueryTokenArchiveCoordinator{backend: b} }
func (c *QueryTokenArchiveCoordinator) Archive(payload []byte) { c.backend.Put(payload); c.pending = payload }
func (c *QueryTokenArchiveCoordinator) Export() []byte { return append([]byte(nil), c.pending...) }
