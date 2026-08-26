package service

import "log-aggregation/internal/store"

type SourceScanResultStreamCoordinator struct { backend *store.SourceScanResultStreamStore }
func NewSourceScanResultStreamCoordinator(b *store.SourceScanResultStreamStore) *SourceScanResultStreamCoordinator { return &SourceScanResultStreamCoordinator{backend: b} }
func (c *SourceScanResultStreamCoordinator) Collect(fail bool) (values []string, err error) {
    results, errs := c.backend.Stream(fail)
    for {
        select {
        case e, ok := <-errs:
            // 失败到达：退出收集并结束这次结果流。
            // errs 关闭（ok==false）表示结果流正常结束、无错误。
            if !ok { return values, nil }
            return values, e
        case value, ok := <-results:
            if !ok { // 结果流结束，等待错误通道关闭以确认是否有失败
                if e, eok := <-errs; eok && e != nil { return values, e }
                return values, nil
            }
            values = append(values, value)
        }
    }
}
