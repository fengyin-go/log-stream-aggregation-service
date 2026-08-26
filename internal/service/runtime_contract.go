package service

import (
    "errors"
    "log-aggregation/internal/store"
)

type RuleEvaluationAssemblyCoordinator struct { backend *store.RuleEvaluationAssemblyStore }
func NewRuleEvaluationAssemblyCoordinator(b *store.RuleEvaluationAssemblyStore) *RuleEvaluationAssemblyCoordinator { return &RuleEvaluationAssemblyCoordinator{backend: b} }
func (c *RuleEvaluationAssemblyCoordinator) Build(key string, fail bool) (item *store.Assembly, err error) { defer func() { if recover() != nil { item, _ = c.backend.Get(key); err = errors.New("assembly failed") } }(); item = c.backend.Build(key, fail); return item, nil }
