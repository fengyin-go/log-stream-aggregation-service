package service

import (
	"errors"
	"log-aggregation/internal/store"
)

type RuleEvaluationAssemblyCoordinator struct {
	backend *store.RuleEvaluationAssemblyStore
}

func NewRuleEvaluationAssemblyCoordinator(b *store.RuleEvaluationAssemblyStore) *RuleEvaluationAssemblyCoordinator {
	return &RuleEvaluationAssemblyCoordinator{backend: b}
}
func (c *RuleEvaluationAssemblyCoordinator) Build(key string, fail bool) (item *store.Assembly, err error) {
	defer func() {
		if recover() != nil {
			// 失败构建不应返回缓存里可能不完整的对象，也不应让其对其他请求可见。
			item = nil
			err = errors.New("assembly failed")
		}
	}()
	item = c.backend.Build(key, fail)
	return item, nil
}
