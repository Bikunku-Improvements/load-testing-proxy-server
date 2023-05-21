package integration

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"google.golang.org/grpc/status"
	"log"
	"sync"
)

type ErrorOccur struct {
	errors map[string]int
	sync   sync.Mutex
}

func (e *ErrorOccur) HandleError(err error) {
	e.sync.Lock()
	defer e.sync.Unlock()

	if errGRPC, ok := status.FromError(err); ok {
		if _, ok := e.errors[errGRPC.Code().String()]; ok {
			e.errors[errGRPC.Code().String()] += 1
		} else {
			e.errors[errGRPC.Code().String()] = 1
		}
	} else {
		errType := fmt.Sprintf("%T", err)
		log.Printf("error type %s; error detail: %s", errType, err.Error())
		if _, ok := e.errors[errType]; ok {
			e.errors[errType] += 1
		} else {
			e.errors[errType] = 1
		}
	}

}

type Throughput struct {
	sizeData  int
	totalData int
	sync      sync.Mutex
}

func (t *Throughput) AddSizeData(v interface{}) error {
	b := new(bytes.Buffer)
	if err := gob.NewEncoder(b).Encode(v); err != nil {
		return err
	}

	t.sync.Lock()
	defer t.sync.Unlock()
	t.sizeData += b.Len()
	t.totalData++

	return nil
}
