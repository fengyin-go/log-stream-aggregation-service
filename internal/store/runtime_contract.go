package store

import "sync"

// SourcePauseStateComplete 是日志源暂停流程的终态。进入终态后，
// 任何较晚到达的异步更新都属于过期更新，必须被拒绝——终态不能后退。
const SourcePauseStateComplete = "complete"

type VersionedState struct {
	Version int
	State   string
}

type SourcePauseVersionStore struct {
	mu     sync.Mutex
	states map[string]VersionedState
}

func NewSourcePauseVersionStore() *SourcePauseVersionStore {
	return &SourcePauseVersionStore{states: make(map[string]VersionedState)}
}

// Update 用带版本的更新推进指定 key 的状态。
//
// 终态不可后退：一旦某个 key 进入 complete 终态，后续到达的任何更新
// （无论版本号新旧）都必须被拒绝，以防止一个较早的异步更新把来源
// 从终态退回活动状态并重复发出状态通知。
//
// 次之，版本号严格单调递增：版本号不大于当前版本的更新视为过期更新，
// 一律拒绝。
//
// 当更新被接受时返回 true，被拒绝（终态后退或过期）时返回 false。
func (s *SourcePauseVersionStore) Update(key string, version int, state string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	current, ok := s.states[key]

	// 终态不能后退：已进入 complete 的来源拒绝一切后续更新。
	if ok && current.State == SourcePauseStateComplete {
		return false
	}

	// 版本必须严格递增：等于或更旧的版本都是过期更新，拒绝。
	if ok && version <= current.Version {
		return false
	}

	s.states[key] = VersionedState{Version: version, State: state}
	return true
}

func (s *SourcePauseVersionStore) Get(key string) VersionedState {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.states[key]
}
