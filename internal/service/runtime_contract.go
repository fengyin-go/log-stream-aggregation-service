package service

import "log-aggregation/internal/store"

type AlertContextArchiveCoordinator struct { backend *store.AlertContextArchiveStore; pending []byte }
func NewAlertContextArchiveCoordinator(b *store.AlertContextArchiveStore) *AlertContextArchiveCoordinator { return &AlertContextArchiveCoordinator{backend: b} }
func (c *AlertContextArchiveCoordinator) Archive(payload []byte) { c.backend.Put(payload); c.pending = payload }
func (c *AlertContextArchiveCoordinator) Export() []byte { return append([]byte(nil), c.pending...) }
