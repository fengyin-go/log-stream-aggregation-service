package service

import "log-aggregation/internal/store"

type QueryTokenArchiveCoordinator struct {
	backend *store.QueryTokenArchiveStore
	pending []byte
}

func NewQueryTokenArchiveCoordinator(b *store.QueryTokenArchiveStore) *QueryTokenArchiveCoordinator {
	return &QueryTokenArchiveCoordinator{backend: b}
}

// Archive 防御性拷贝 payload：保存后调用方复用同一输入切片写入新令牌，
// 也不会改写已归档的 pending 与后端缓存。
func (c *QueryTokenArchiveCoordinator) Archive(payload []byte) {
	c.backend.Put(payload)
	c.pending = append([]byte(nil), payload...)
}
func (c *QueryTokenArchiveCoordinator) Export() []byte { return append([]byte(nil), c.pending...) }
