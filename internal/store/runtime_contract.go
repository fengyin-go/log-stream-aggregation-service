package store

import "sync"

type Assembly struct {
	ID    string
	Parts map[string]string
}
type RuleEvaluationAssemblyStore struct {
	mu    sync.Mutex
	items map[string]*Assembly
}

func NewRuleEvaluationAssemblyStore() *RuleEvaluationAssemblyStore {
	return &RuleEvaluationAssemblyStore{items: make(map[string]*Assembly)}
}
func (s *RuleEvaluationAssemblyStore) Build(key string, fail bool) (assembled *Assembly) {
	assembled = &Assembly{ID: key}
	if fail {
		panic("payload decode")
	}
	assembled.Parts = map[string]string{"header": "ready"}
	// 只有完整对象才落缓存：失败构建（panic）发生在写缓存之前，
	// 不让只有编号的不完整对象对其他请求可见。
	s.mu.Lock()
	s.items[key] = assembled
	s.mu.Unlock()
	return assembled
}
func (s *RuleEvaluationAssemblyStore) Get(key string) (*Assembly, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.items[key]
	return item, ok
}
func (s *RuleEvaluationAssemblyStore) Delete(key string) {
	s.mu.Lock()
	delete(s.items, key)
	s.mu.Unlock()
}
