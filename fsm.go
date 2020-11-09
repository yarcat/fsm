package fsm

import "log"

// FSM is a finite-state machine. It implements state transitions and events.
type FSM struct {
	current StateType

	transitions    Transitions
	states         States
	defaultHandler State
}

// New returns new configured finite-state machine.
func New(startState StateType, transitions Transitions, states States, defaultHandler State) *FSM {
	if defaultHandler == nil {
		defaultHandler = DefaultHandler
	}
	fsm := &FSM{
		transitions:    transitions,
		states:         states,
		defaultHandler: defaultHandler,
	}
	fsm.handler(startState).Enter()
	fsm.current = startState
	log.Printf("fsm %p: initialized: state %v", fsm, startState)
	return fsm
}

// Send sends the event to the finite-state machine. The machine changes its
// state accordingly to the transitions table.
func (fsm *FSM) Send(event EventType) {
	log.Printf("fsm %p: state %v: received %v ", fsm, fsm.current, event)
	next, ok := fsm.transitions[StateEvent{State: fsm.current, Event: event}]
	if !ok {
		return
	}
	fsm.change(next)
}

// change changes the state calling Leave and Enter state methods. The method
// does not fire the state handler if next and current states are the same.
func (fsm *FSM) change(next StateType) {
	if fsm.current == next {
		return
	}
	log.Printf("fsm %p: state %v: next state %v", fsm, fsm.current, next)
	fsm.handler(fsm.current).Leave()
	fsm.handler(next).Enter()
	fsm.current = next
}

func (fsm *FSM) handler(state StateType) State {
	if h, ok := fsm.states[state]; ok {
		return h
	}
	return fsm.defaultHandler
}
