package load_test

import (
	"sort"
	"sync"
)

type EndToEndResponseTime struct {
	times  []float64
	sync   sync.Mutex
	newLoc map[int]int
}

type Throughput struct {
	count int
	size  int
	sync  sync.Mutex
}

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

func median(data []float64) float64 {
	dataCopy := make([]float64, len(data))
	copy(dataCopy, data)

	sort.Float64s(dataCopy)

	var median float64
	l := len(dataCopy)
	if l == 0 {
		return 0
	} else if l%2 == 0 {
		median = (dataCopy[l/2-1] + dataCopy[l/2]) / 2
	} else {
		median = dataCopy[l/2]
	}

	return median
}
