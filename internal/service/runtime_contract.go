package service

import (
	"fmt"
	"sync"

	"log-aggregation/internal/store"
)

// SourcePauseVersionCoordinator 协调日志源暂停流程的带版本状态更新。
//
// 异步更新可能乱序到达：一个较早的更新可能在来源已经进入暂停终态之后
// 才到达。coordinator 依赖后端 store 拒绝过期更新与终态后退，自身则保证
// 只有被后端接受的更新才会产生副作用（记录 effect、发出状态通知），
// 从而避免把来源从终态退回活动状态并重复发出状态通知。
type SourcePauseVersionCoordinator struct {
	backend *store.SourcePauseVersionStore
	mu      sync.Mutex
	effects map[string]bool
}

func NewSourcePauseVersionCoordinator(b *store.SourcePauseVersionStore) *SourcePauseVersionCoordinator {
	return &SourcePauseVersionCoordinator{backend: b, effects: make(map[string]bool)}
}

// apply 提交一次带版本的更新。仅当后端接受该更新（既非过期、也未触发
// 终态后退）时才记录 effect 等副作用；被拒绝的过期更新不产生任何影响。
func (c *SourcePauseVersionCoordinator) apply(key string, version int, state string) bool {
	if !c.backend.Update(key, version, state) {
		// 过期更新或终态后退，被拒绝：不发出状态通知、不记录 effect。
		return false
	}
	c.mu.Lock()
	c.effects[fmt.Sprintf("%s-%d", key, version)] = true
	c.mu.Unlock()
	return true
}

// CompleteThenLate 模拟乱序到达：先令来源进入 complete 终态（版本 2），
// 随后一个较早的 running 更新（版本 1）迟到。终态不能后退，因此第二次
// 更新必须被拒绝——来源停留在 complete，且只会产生一次 effect。
func (c *SourcePauseVersionCoordinator) CompleteThenLate(key string) {
	c.apply(key, 2, store.SourcePauseStateComplete)
	c.apply(key, 1, "running")
}

func (c *SourcePauseVersionCoordinator) EffectCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.effects)
}
