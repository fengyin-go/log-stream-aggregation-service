package service

import (
    "errors"
    "log-aggregation/internal/store"
)

type LogMetadataAssemblyCoordinator struct { backend *store.LogMetadataAssemblyStore }
func NewLogMetadataAssemblyCoordinator(b *store.LogMetadataAssemblyStore) *LogMetadataAssemblyCoordinator { return &LogMetadataAssemblyCoordinator{backend: b} }
// Build 组装日志元数据。解析失败时恢复成普通错误返回，但失败对象绝不外泄：
// 即便后端缓存里意外残留了半初始化对象，这里也不会把它取出来交给调用方。
func (c *LogMetadataAssemblyCoordinator) Build(key string, fail bool) (item *store.Assembly, err error) { defer func() { if r := recover(); r != nil { c.backend.Delete(key); item = nil; err = errors.New("assembly failed") } }(); item = c.backend.Build(key, fail); return item, nil }
