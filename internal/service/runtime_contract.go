package service

import "log-aggregation/internal/store"

// IngestRateSnapshotCoordinator 在新日志写入与摄取速率统计同时发生时，
// 负责取得一个稳定的摄取速率快照供上层读取。
type IngestRateSnapshotCoordinator struct {
	backend *store.IngestRateSnapshotStore
}

func NewIngestRateSnapshotCoordinator(b *store.IngestRateSnapshotStore) *IngestRateSnapshotCoordinator {
	return &IngestRateSnapshotCoordinator{backend: b}
}

// CountDuringUpdate 返回采集开始那一刻的摄取速率计数。
//
// 摄取速率统计（取快照）和新日志写入（Increment）可能并发发生。
// 快照本身是一个冻结的值副本，因此即便随后有新日志写入，上层
// 读取期间报出的数字也始终是采集开始时的值，不会变化。
func (c *IngestRateSnapshotCoordinator) CountDuringUpdate() int {
	return c.backend.Snapshot()
}
