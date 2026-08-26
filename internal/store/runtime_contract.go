package store

import "sync"

type Assembly struct { ID string; Parts map[string]string }
type LogMetadataAssemblyStore struct { mu sync.Mutex; items map[string]*Assembly }
func NewLogMetadataAssemblyStore() *LogMetadataAssemblyStore { return &LogMetadataAssemblyStore{items: make(map[string]*Assembly)} }
// Build 组装日志元数据。候选对象只有在组装成功后才进入可见缓存，
// 解析失败（panic）发生在发布之前，因此半初始化对象绝不会留在可共享缓存里。
func (s *LogMetadataAssemblyStore) Build(key string, fail bool) *Assembly { candidate := &Assembly{ID: key}; if fail { panic("payload decode") }; candidate.Parts = map[string]string{"header": "ready"}; s.mu.Lock(); s.items[key] = candidate; s.mu.Unlock(); return candidate }
func (s *LogMetadataAssemblyStore) Get(key string) (*Assembly, bool) { s.mu.Lock(); defer s.mu.Unlock(); item, ok := s.items[key]; return item, ok }
func (s *LogMetadataAssemblyStore) Delete(key string) { s.mu.Lock(); delete(s.items, key); s.mu.Unlock() }
