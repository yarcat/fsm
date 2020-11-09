package fsm

import (
	"reflect"
	"testing"
)

func TestFSMSend(t *testing.T) {
	transitions := Transitions{
		{"st", "ev"}: "st2",
	}
	for _, tc := range []struct {
		name                   string
		fsm                    *FSM
		evt                    EventType
		wantState              StateType
		wantEnters, wantLeaves int
	}{
		{"nil transitions", New("st", nil, nil, new(countingHandler)), "ev", "st", 0, 0},
		{"wrong state", New("wrong", transitions, nil, new(countingHandler)), "ev", "wrong", 0, 0},
		{"unexpected event", New("st", transitions, nil, new(countingHandler)), "evUnexpected", "st", 0, 0},
		{"changes", New("st", transitions, nil, new(countingHandler)), "ev", "st2", 1, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.fsm.Send(tc.evt)
			state := tc.fsm.current
			if state != tc.wantState {
				t.Errorf("Send(%v); fsm.current = %v, want = %v",
					tc.evt, tc.fsm.current, tc.wantState)
			}
			enters, leaves := tc.fsm.defaultHandler.(*countingHandler).counters()
			// Compensate enters for the call from New().
			enters--
			if enters != tc.wantEnters || leaves != tc.wantLeaves {
				t.Errorf("Send(%v); got (enters,leaves) = (%v,%v), want = (%v,%v)",
					tc.evt, enters, leaves, tc.wantEnters, tc.wantLeaves)
			}
		})
	}
}

type testHandler struct{}

func (testHandler) Enter() {}
func (testHandler) Leave() {}

func TestNew(t *testing.T) {
	var (
		testHandler testHandler
		withState   countingHandler
	)
	for _, tc := range []struct {
		name                               string
		startState                         StateType
		transitions                        Transitions
		states                             States
		defaultHandler, wantDefaultHandler State
		wantEnters, wantLeaves             int
	}{
		{"nil", "", nil, nil, nil, DefaultHandler, 0, 0},
		{"default handler", "", nil, nil, DefaultHandler, DefaultHandler, 0, 0},
		{"custom handler", "", nil, nil, testHandler, testHandler, 0, 0},
		{"states", "", nil, States{}, nil, DefaultHandler, 0, 0},
		{"transitions", "", Transitions{}, nil, nil, DefaultHandler, 0, 0},
		{"start state", "start", nil, nil, nil, DefaultHandler, 0, 0},
		{"all", "start", Transitions{}, States{}, testHandler, testHandler, 0, 0},
		{"calls enter", "start", Transitions{}, States{}, &withState, &withState, 1, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fsm := New(tc.startState, tc.transitions, tc.states, tc.defaultHandler)
			if fsm.current != tc.startState {
				t.Errorf("New().current = %v, want = %v", fsm.current, tc.startState)
			}
			if !reflect.DeepEqual(fsm.transitions, tc.transitions) {
				t.Errorf("New().transitions = %v, want = %v", fsm.transitions, tc.transitions)
			}
			if !reflect.DeepEqual(fsm.states, tc.states) {
				t.Errorf("New().states = %v, want = %v", fsm.states, tc.states)
			}
			if fsm.defaultHandler != tc.wantDefaultHandler {
				t.Errorf("New().defaultHandler = %v, want = %v", fsm.defaultHandler, tc.defaultHandler)
			}
			if tc.defaultHandler != &withState {
				return
			}
			enters, leaves := withState.counters()
			if enters != tc.wantEnters || leaves != tc.wantLeaves {
				t.Errorf("New() count(Enter, Leave) = (%v,%v), want = (%v,%v)",
					enters, leaves, tc.wantEnters, tc.wantLeaves)
			}
		})
	}
}
