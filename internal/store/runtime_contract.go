package store

import "sync"

type VersionedState struct {
	Version int
	State   string
}

// terminalStates 标记告警解决的终态。一旦进入终态，迟到的旧回调不得将其覆盖回处理中。
var terminalStates = map[string]bool{"complete": true}

func isTerminalState(state string) bool { return terminalStates[state] }

type AlertResolutionVersionStore struct {
	mu     sync.Mutex
	states map[string]VersionedState
}

func NewAlertResolutionVersionStore() *AlertResolutionVersionStore {
	return &AlertResolutionVersionStore{states: make(map[string]VersionedState)}
}

// Update 以版本号守护状态写入：旧版本不能覆盖新版本，终态不能被非终态回调回退。
// 仅在实际写入成功时返回 true，供上层据此决定是否触发通知副作用。
func (s *AlertResolutionVersionStore) Update(key string, version int, state string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	cur, ok := s.states[key]
	if ok {
		// 旧版本不能覆盖终态：版本号更旧的迟到回调直接丢弃。
		if version < cur.Version {
			return false
		}
		// 终态不可回退：已为终态时，非终态回调一律不覆盖。
		if isTerminalState(cur.State) && !isTerminalState(state) {
			return false
		}
	}
	s.states[key] = VersionedState{Version: version, State: state}
	return true
}

func (s *AlertResolutionVersionStore) Get(key string) VersionedState {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.states[key]
}
