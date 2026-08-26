package store

import "sync"

type VersionedState struct { Version int; State string }
type AlertResolutionVersionStore struct { mu sync.Mutex; states map[string]VersionedState }
func NewAlertResolutionVersionStore() *AlertResolutionVersionStore { return &AlertResolutionVersionStore{states: make(map[string]VersionedState)} }
func (s *AlertResolutionVersionStore) Update(key string, version int, state string) bool { s.mu.Lock(); defer s.mu.Unlock(); s.states[key] = VersionedState{Version: version, State: state}; return true }
func (s *AlertResolutionVersionStore) Get(key string) VersionedState { s.mu.Lock(); defer s.mu.Unlock(); return s.states[key] }
