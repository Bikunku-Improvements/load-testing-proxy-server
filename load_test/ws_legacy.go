package load_test

import (
	"encoding/json"
	"fmt"
	websocket_dialler "github.com/fasthttp/websocket"
	"load-testing-proxy-server/entity"
	"log"
	"os"
	"sync"
	"time"
	"unsafe"
)

func WSLegacyClientTest(concurrentUser int, receiveMessagePerClient int) {
	log.Printf("starting ws legacy load test with %d concurrent user and %d receive message per client", concurrentUser, receiveMessagePerClient)
	var wg sync.WaitGroup

	avgTime := EndToEndResponseTime{
		sync:   sync.Mutex{},
		newLoc: make(map[int]int),
	}

	errorsOccur := ErrorOccur{
		errors: make(map[string]int),
		sync:   sync.Mutex{},
	}

	for i := 0; i < concurrentUser; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			addr := os.Getenv("LEGACY_BIKUNKU_SERVER")
			dial, _, err := websocket_dialler.DefaultDialer.Dial(fmt.Sprintf("ws://%s/bus/stream", addr), nil)
			if err != nil {
				errorsOccur.HandleError(err)
				return
			}
			defer dial.Close()

			var totalRequest int
			for totalRequest < receiveMessagePerClient {
				_, msg, err := dial.ReadMessage()
				if err != nil {
					errorsOccur.HandleError(err)
					return
				}

				var location []entity.BusLocationWSLegacyResponse
				err = json.Unmarshal(msg, &location)
				if err != nil {
					errorsOccur.HandleError(err)
					return
				}

				for _, v := range location {
					avgTime.sync.Lock()
					if _, ok := avgTime.newLoc[v.Id]; !ok {
						responseTime := time.Since(v.Timestamp)
						avgTime.newLoc[v.Id] = v.Id
						avgTime.times = append(avgTime.times, responseTime.Seconds())
						totalRequest++

						log.Printf("message received from websocket legacy with id=%d with latency: %s", v.Id, responseTime)
					}
					avgTime.sync.Unlock()
				}
			}
		}()
	}
	wg.Wait()

	log.Printf("median response time: %v second", median(avgTime.times))
	log.Printf("error occured: %v", errorsOccur.errors)
}

// Deprecated: please use WSLegacyDriverTest in integration package
func WSLegacyDriverTest(concurrentUser, sendMessagePerClient int) {
	log.Printf("starting ws legacy driver load test with %d concurrent user  and %d send message per client", concurrentUser, sendMessagePerClient)

	throughput := Throughput{
		sync: sync.Mutex{},
	}

	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < concurrentUser; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			addr := os.Getenv("LEGACY_BIKUNKU_SERVER")
			token := os.Getenv("LEGACY_BIKUNKU_TOKEN")

			url := fmt.Sprintf("ws://%s/bus/stream?type=%s&token=%s", addr, "driver", token)
			dial, _, err := websocket_dialler.DefaultDialer.Dial(url, nil)
			if err != nil {
				log.Println(err)
				return
			}
			defer dial.Close()

			var totalRequest int
			for totalRequest < sendMessagePerClient {
				data := []byte(`{"long": 91, "lat": 20, "speed": 0.00, "heading": 0}`)
				err = dial.WriteMessage(1, data)
				if err != nil {
					log.Println(err)
					return
				}

				throughput.sync.Lock()
				throughput.count++
				throughput.size += int(unsafe.Sizeof(data))
				throughput.sync.Unlock()
				totalRequest++
			}
		}()
	}
	wg.Wait()
	log.Printf("total size of sent data: %v kb/s", float64(throughput.size)/1000/time.Since(start).Seconds())
}
