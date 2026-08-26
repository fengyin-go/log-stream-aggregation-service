package service

import (
	"fmt"
	"sync"

	"log-aggregation/internal/store"
)

// AlertResolutionVersionCoordinator 协调告警解决过程中的版本化状态写入。
// 关键约束：解决终态一旦确立，迟到的旧回调不得将其覆盖回处理中，且解决通知只能发送一次。
type AlertResolutionVersionCoordinator struct {
	backend *store.AlertResolutionVersionStore
	mu      sync.Mutex
	effects map[string]bool
}

func NewAlertResolutionVersionCoordinator(b *store.AlertResolutionVersionStore) *AlertResolutionVersionCoordinator {
	return &AlertResolutionVersionCoordinator{backend: b, effects: make(map[string]bool)}
}

// apply 写入一次状态更新，仅在实际写入成功时记录副作用（对应通知动作）。
// 旧版本或回退终态的迟到回调会被 backend 拒绝，从而既不改变终态也不重复触发通知。
func (c *AlertResolutionVersionCoordinator) apply(key string, version int, state string) {
	if !c.backend.Update(key, version, state) {
		return
	}
	c.mu.Lock()
	c.effects[fmt.Sprintf("%s-%d", key, version)] = true
	c.mu.Unlock()
}

// CompleteThenLate 模拟解决流程：终态（complete）先落地，随后迟到的旧回调（running）到达。
// 修复后迟到回调无法覆盖终态，通知副作用只记录一次。
func (c *AlertResolutionVersionCoordinator) CompleteThenLate(key string) {
	c.apply(key, 2, "complete")
	c.apply(key, 1, "running")
}

// EffectCount 返回已实际生效的写入次数，即通知动作的次数。
func (c *AlertResolutionVersionCoordinator) EffectCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.effects)
}
