package service

import "log-aggregation/internal/store"

type AlertContextArchiveCoordinator struct { backend *store.AlertContextArchiveStore; pending []byte }
func NewAlertContextArchiveCoordinator(b *store.AlertContextArchiveStore) *AlertContextArchiveCoordinator { return &AlertContextArchiveCoordinator{backend: b} }
// Archive 将告警上下文写入归档。调用方可能复用入参所在的原始切片，
// 因此归档前必须把上下文从输入缓冲区彻底分离：交给 backend 持有
// 独立拷贝，pending 也只保留自己的拷贝，避免与导出/缓存中的旧上下文
// 共享底层数组而相互污染。
func (c *AlertContextArchiveCoordinator) Archive(payload []byte) {
	owned := append([]byte(nil), payload...)
	c.backend.Put(owned)
	c.pending = owned
}

// Export 返回已存档上下文的独立拷贝，读取方修改返回值不会影响缓存。
func (c *AlertContextArchiveCoordinator) Export() []byte {
	return append([]byte(nil), c.pending...)
}
