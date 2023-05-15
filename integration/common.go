package integration

import (
	"fmt"
	"sync"
)

type ErrorOccur struct {
	errors map[string]int
	sync   sync.Mutex
}

func (e *ErrorOccur) HandleError(err error) {
	e.sync.Lock()
	defer e.sync.Unlock()

	errType := fmt.Sprintf("%T", err)
	if _, ok := e.errors[errType]; ok {
		e.errors[errType] += 1
	} else {
		e.errors[errType] = 1
	}
}
