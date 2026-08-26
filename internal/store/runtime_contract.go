package store

import "sync"

type Assembly struct { ID string; Parts map[string]string }
type ParsedEntryAssemblyStore struct { mu sync.Mutex; items map[string]*Assembly }
func NewParsedEntryAssemblyStore() *ParsedEntryAssemblyStore { return &ParsedEntryAssemblyStore{items: make(map[string]*Assembly)} }
func (s *ParsedEntryAssemblyStore) Build(key string, fail bool) *Assembly {
	// 先在缓存外完成组装，确保只有完整条目才会被提交。
	// 若中途失败（panic），半成品不会进入缓存，后续 Get 也就读不到字段不全的条目。
	candidate := &Assembly{Parts: map[string]string{"header": "ready"}}
	if fail {
		panic("payload decode")
	}
	candidate.ID = key
	s.mu.Lock()
	s.items[key] = candidate
	s.mu.Unlock()
	return candidate
}
func (s *ParsedEntryAssemblyStore) Get(key string) (*Assembly, bool) { s.mu.Lock(); defer s.mu.Unlock(); item, ok := s.items[key]; return item, ok }
func (s *ParsedEntryAssemblyStore) Delete(key string) { s.mu.Lock(); delete(s.items, key); s.mu.Unlock() }
