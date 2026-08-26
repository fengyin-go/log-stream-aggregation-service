package store

import "errors"

type SourceScanResultStreamStore struct{}
func NewSourceScanResultStreamStore() *SourceScanResultStreamStore { return &SourceScanResultStreamStore{} }
func (s *SourceScanResultStreamStore) Stream(fail bool) (<-chan string, <-chan error) {
    results := make(chan string); errs := make(chan error, 1)
    go func() {
        defer close(results)
        defer close(errs)
        if fail { errs <- errors.New("partition unavailable"); return }
        results <- "ready"
    }()
    return results, errs
}
