package store

import "sync"

type VersionedState struct { Version int; State string }
type NotificationReadVersionStore struct { mu sync.Mutex; states map[string]VersionedState }
func NewNotificationReadVersionStore() *NotificationReadVersionStore { return &NotificationReadVersionStore{states: make(map[string]VersionedState)} }
// Update 仅接受严格递增的版本，拒绝版本倒退或重复回调，
// 从而保护已写入的已读终态不被首轮延迟回调覆盖。
func (s *NotificationReadVersionStore) Update(key string, version int, state string) bool { s.mu.Lock(); defer s.mu.Unlock(); if cur, ok := s.states[key]; ok && version <= cur.Version { return false }; s.states[key] = VersionedState{Version: version, State: state}; return true }
func (s *NotificationReadVersionStore) Get(key string) VersionedState { s.mu.Lock(); defer s.mu.Unlock(); return s.states[key] }
