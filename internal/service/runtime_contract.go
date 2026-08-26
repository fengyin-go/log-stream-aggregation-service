package service

import (
    "sync"
    "fmt"
    "log-aggregation/internal/store"
)

type NotificationReadVersionCoordinator struct { backend *store.NotificationReadVersionStore; mu sync.Mutex; effects map[string]bool }
func NewNotificationReadVersionCoordinator(b *store.NotificationReadVersionStore) *NotificationReadVersionCoordinator { return &NotificationReadVersionCoordinator{backend: b, effects: make(map[string]bool)} }
// apply 仅在后端真正接受了该版本（非倒退、非重复）时才记录确认副作用，
// 保证终态不会被首轮延迟回调覆盖，且外部确认动作只生效一次。
func (c *NotificationReadVersionCoordinator) apply(key string, version int, state string) { if !c.backend.Update(key, version, state) { return }; c.mu.Lock(); c.effects[fmt.Sprintf("%s-%d", key, version)] = true; c.mu.Unlock() }
func (c *NotificationReadVersionCoordinator) CompleteThenLate(key string) { c.apply(key, 2, "complete"); c.apply(key, 1, "running") }
func (c *NotificationReadVersionCoordinator) EffectCount() int { c.mu.Lock(); defer c.mu.Unlock(); return len(c.effects) }
