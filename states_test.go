package fsm

import (
	"testing"
)

func TestDefault(t *testing.T) {
	var stateHandler State = DefaultHandler
	stateHandler.Enter()
	stateHandler.Leave()
}

func TestComposite(t *testing.T) {
	for _, tc := range []struct {
		name                   string
		state                  Composite
		wantEnters, wantLeaves int
	}{
		{"nil", nil, 0, 0},
		{"empty", counters(0), 0, 0},
		{"single", counters(1), 1, 1},
		{"many", counters(100), 100, 100},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.state.Enter()
			tc.state.Leave()
			enters, leaves := collect(tc.state)
			if enters != tc.wantEnters || leaves != tc.wantLeaves {
				t.Errorf("got (enters,leaves) = (%v,%v), want (%v,%v)",
					enters, leaves, tc.wantEnters, tc.wantLeaves)
			}
		})
	}
}

func TestCompose(t *testing.T) {
	for _, tc := range []struct {
		name    string
		states  []State
		wantLen int
	}{
		{"nil", nil, 0},
		{"single", []State{DefaultHandler}, 1},
		{"few", []State{DefaultHandler, DefaultHandler, DefaultHandler}, 3},
	} {
		t.Run(tc.name, func(t *testing.T) {
			st := Compose(tc.states...)
			if len(st) != tc.wantLen {
				t.Errorf("len(st) = %v, want = %v", len(st), tc.wantLen)
			}
		})
	}
}

func counters(n int) (s Composite) {
	for i := 0; i < n; i++ {
		s = append(s, &countingHandler{})
	}
	return
}

func collect(s Composite) (enters, leaves int) {
	for _, st := range s {
		e, l := st.(*countingHandler).counters()
		enters += e
		leaves += l
	}
	return
}

func TestMachineProvider(t *testing.T) {
	const (
		stInit        StateType = "stInit"
		stStart       StateType = "stStart"
		evInitialized EventType = "evInitialized"
	)
	transitions := Transitions{
		When(stInit, evInitialized): stStart,
	}
	provider := new(MachineProvider)
	fsm := New(stInit, transitions, nil, nil)
	provider.Set(fsm)
	if fsm.current != stInit {
		t.Errorf("current = %v, want = %v", fsm.current, stInit)
	}
	provider.Send(evInitialized)
	if fsm.current != stStart {
		t.Errorf("current = %v, want = %v", fsm.current, stStart)
	}
}

func TestExpiring(t *testing.T) {
	// Setup test "call after" implementation.
	var (
		cancelCalled bool
		scheduledFn  func()
	)
	testAfter := func(f func()) func() {
		scheduledFn = f
		return func() {
			cancelCalled = true
		}
	}

	const (
		stInit        StateType = "stInit"
		stWaitTimeout StateType = "stWaitTimeout"
		stFinal       StateType = "stFinal"
		evInitialized EventType = "evInitialized"
		evTimeout     EventType = "evTimeout"
	)

	transitions := Transitions{
		When(stInit, evInitialized):    stWaitTimeout,
		When(stWaitTimeout, evTimeout): stFinal,
	}

	provider := new(MachineProvider)
	states := States{
		stWaitTimeout: NewExpiring(provider, testAfter, evTimeout),
	}
	fsm := New(stInit, transitions, states, nil)
	provider.Set(fsm)
	fsm.Send(evInitialized)

	if fsm.current != stWaitTimeout {
		t.Errorf("fsm.current = %v, want = %v", fsm.current, evInitialized)
	}
	if cancelCalled {
		t.Errorf("cancel called = %v, want = %v", cancelCalled, false)
	}
	scheduledFn()
	if fsm.current != stFinal {
		t.Errorf("fsm.current = %v, want = %v", fsm.current, stFinal)
	}
	if !cancelCalled {
		t.Errorf("cancel called = %v, want = %v", cancelCalled, true)
	}
}
