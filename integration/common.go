package integration

import "sync"

type ErrorOccur struct {
	errors map[string]int
	sync   sync.Mutex
}

func (e *ErrorOccur) HandleError(err error) {
	e.sync.Lock()
	defer e.sync.Unlock()

	if _, ok := e.errors[err.Error()]; ok {
		e.errors[err.Error()] += 1
	} else {
		e.errors[err.Error()] = 1
	}
}
