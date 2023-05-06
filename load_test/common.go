package load_test

import (
	"sort"
	"sync"
)

type EndToEndResponseTime struct {
	times []float64
	sync  sync.Mutex
}

type Throughput struct {
	count int
	size  int
	sync  sync.Mutex
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
