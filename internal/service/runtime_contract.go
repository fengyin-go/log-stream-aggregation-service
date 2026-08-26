package service

import "log-aggregation/internal/store"

type SourceCounterSnapshotCoordinator struct { backend *store.SourceCounterSnapshotStore }
func NewSourceCounterSnapshotCoordinator(b *store.SourceCounterSnapshotStore) *SourceCounterSnapshotCoordinator { return &SourceCounterSnapshotCoordinator{backend: b} }
// CountDuringUpdate 模拟一次采集轮次：先固定采集开始时的来源计数快照，
// 采集期间新写入的日志通过 Increment() 更新实时计数，但只影响下一轮快照。
// 因此必须先取快照，再触发本轮的增量，保证本轮交给上层的是开始时的计数。
func (c *SourceCounterSnapshotCoordinator) CountDuringUpdate() int {
	snapshot := c.backend.Snapshot()
	done := make(chan struct{})
	go func() { c.backend.Increment(); close(done) }()
	<-done
	return *snapshot
}
