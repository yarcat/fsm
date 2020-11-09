package fsm

type countingHandler struct{ enters, leaves int }

func (h *countingHandler) Enter()               { h.enters++ }
func (h *countingHandler) Leave()               { h.leaves++ }
func (h *countingHandler) counters() (int, int) { return h.enters, h.leaves }
