package fsm

type (
	// EventType is a type of FSM events.
	EventType string

	// Transitions describes state transitions in a form of a mapping from
	// a state and an event to a new state.
	Transitions map[StateEvent]StateType
	// StateEvent is a transitions key in the form of a state and an event.
	StateEvent struct {
		State StateType
		Event EventType
	}
)

// When is a helper function with an intention to make transition map creation
// more readable.
//
// Example:
//
//	const (
//		stInit        StateType = "stInit"
//		stStart       StateType = "stStart"
//		evInitialized EventType = "evInitialized"
//	)
//	transitions := Transitions{
//		When(stInit, evInitialized): stStart,
//	}
func When(st StateType, ev EventType) StateEvent {
	return StateEvent{State: st, Event: ev}
}
