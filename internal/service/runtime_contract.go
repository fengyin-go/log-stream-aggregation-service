package service

import (
    "sync"
    "fmt"
    "log-aggregation/internal/store"
)

type AlertResolutionVersionCoordinator struct { backend *store.AlertResolutionVersionStore; mu sync.Mutex; effects map[string]bool }
func NewAlertResolutionVersionCoordinator(b *store.AlertResolutionVersionStore) *AlertResolutionVersionCoordinator { return &AlertResolutionVersionCoordinator{backend: b, effects: make(map[string]bool)} }
func (c *AlertResolutionVersionCoordinator) apply(key string, version int, state string) { c.backend.Update(key, version, state); c.mu.Lock(); c.effects[fmt.Sprintf("%s-%d", key, version)] = true; c.mu.Unlock() }
func (c *AlertResolutionVersionCoordinator) CompleteThenLate(key string) { c.apply(key, 2, "complete"); c.apply(key, 1, "running") }
func (c *AlertResolutionVersionCoordinator) EffectCount() int { c.mu.Lock(); defer c.mu.Unlock(); return len(c.effects) }
