package service

import (
    "sync"
    "fmt"
    "log-aggregation/internal/store"
)

type SourcePauseVersionCoordinator struct { backend *store.SourcePauseVersionStore; mu sync.Mutex; effects map[string]bool }
func NewSourcePauseVersionCoordinator(b *store.SourcePauseVersionStore) *SourcePauseVersionCoordinator { return &SourcePauseVersionCoordinator{backend: b, effects: make(map[string]bool)} }
func (c *SourcePauseVersionCoordinator) apply(key string, version int, state string) { c.backend.Update(key, version, state); c.mu.Lock(); c.effects[fmt.Sprintf("%s-%d", key, version)] = true; c.mu.Unlock() }
func (c *SourcePauseVersionCoordinator) CompleteThenLate(key string) { c.apply(key, 2, "complete"); c.apply(key, 1, "running") }
func (c *SourcePauseVersionCoordinator) EffectCount() int { c.mu.Lock(); defer c.mu.Unlock(); return len(c.effects) }
