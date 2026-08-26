package store

import "sync"

type VersionedState struct { Version int; State string }
type SourcePauseVersionStore struct { mu sync.Mutex; states map[string]VersionedState }
func NewSourcePauseVersionStore() *SourcePauseVersionStore { return &SourcePauseVersionStore{states: make(map[string]VersionedState)} }
func (s *SourcePauseVersionStore) Update(key string, version int, state string) bool { s.mu.Lock(); defer s.mu.Unlock(); s.states[key] = VersionedState{Version: version, State: state}; return true }
func (s *SourcePauseVersionStore) Get(key string) VersionedState { s.mu.Lock(); defer s.mu.Unlock(); return s.states[key] }
