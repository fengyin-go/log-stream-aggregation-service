package service

import (
    "errors"
    "log-aggregation/internal/store"
)

type ParsedEntryAssemblyCoordinator struct {
    backend *store.ParsedEntryAssemblyStore
}

func NewParsedEntryAssemblyCoordinator(b *store.ParsedEntryAssemblyStore) *ParsedEntryAssemblyCoordinator {
    return &ParsedEntryAssemblyCoordinator{backend: b}
}

// Build 组装一条日志解析条目。组装中途失败时，缓存里残留的半成品会被清掉，
// 调用方只会拿到错误，而不会读到字段不全的条目。
func (c *ParsedEntryAssemblyCoordinator) Build(key string, fail bool) (item *store.Assembly, err error) {
    defer func() {
        if r := recover(); r != nil {
            // 任何失败路径都不该留下半成品：清掉缓存里的残留，避免后续读取看到字段不全的条目。
            c.backend.Delete(key)
            item = nil
            err = errors.New("assembly failed")
        }
    }()
    item = c.backend.Build(key, fail)
    return item, nil
}
