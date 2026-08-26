package store

import "errors"

type ExportChunkResultStreamStore struct{}
func NewExportChunkResultStreamStore() *ExportChunkResultStreamStore { return &ExportChunkResultStreamStore{} }
func (s *ExportChunkResultStreamStore) Stream(fail bool) (<-chan string, <-chan error) {
    results := make(chan string); errs := make(chan error, 1)
    go func() {
        defer close(results)
        defer close(errs)
        if fail { errs <- errors.New("partition unavailable"); return }
        results <- "ready"
    }()
    return results, errs
}
