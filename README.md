# fsm-go
Golang finite-state machine.

# Examples

## Comprehensive

In [this example] we do the following:

* utilize a back-reference to our state maching using `MachineProvider`;
* use `Expiring` state handler to timeout a state;
* create custom state handlers that print the states and terminate the app;
* join the handlers using `Compose`.

[this example]: example_test.go

## Simple

In this example we start with `stInit`, and then send `evInitialized`, which
switches the state into `stFinal`.

```go
const (
	stInit        fsm.StateType = "stInit"
	stFinal       fsm.StateType = "stFinal"
	evInitialized fsm.EventType = "evInitialized"
)

transitions := fsm.Transitions{
	fsm.When(stInit, evInitialized): stFinal,
}

fsm := fsm.New(stInit, transitions, nil, nil)
fsm.Send(evInitialized)
```

## Expiring states

In this example we use an `Expiring` state handler, which sends events to the
FSM with a delay. Since timer events are asynchronous, we recommend to use the
asynchronous FSM implementation.


```golang
const (
	stInit         StateType = "stInit"
	stWaitTimeout  StateType = "stWaitTimeout"
	stFinal        StateType = "stFinal"
	evInitialized  EventType = "evInitialized"
	evTimeout      EventType = "evTimeout"
)

transitions := Transitions{
	When(stInit, evInitialized):    stWaitTimeout,
	When(stWaitTimeout, evTimeout): stFinal,
}

provider := new(MachineProvider)

states := States{
	stWaitTimeout: NewExpiring(provider, After(time.Second), evTimeout),
}

fsm := fsm.NewAsync(stInit, transitions, states, nil)
provider.Set(fsm)
fsm.Send(evInitialized)
fsm.Run()
```
