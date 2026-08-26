package service

import "log-aggregation/internal/store"

type IngestPayloadSnapshotCoordinator struct { backend *store.IngestPayloadSnapshotStore; pending []byte }
func NewIngestPayloadSnapshotCoordinator(b *store.IngestPayloadSnapshotStore) *IngestPayloadSnapshotCoordinator { return &IngestPayloadSnapshotCoordinator{backend: b} }
func (c *IngestPayloadSnapshotCoordinator) Archive(payload []byte) { c.backend.Put(payload); c.pending = append([]byte(nil), payload...) }
func (c *IngestPayloadSnapshotCoordinator) Export() []byte { return append([]byte(nil), c.pending...) }
