package load_test

import (
	"context"
	"fmt"
	websocket_dialler "github.com/fasthttp/websocket"
	"load-testing-proxy-server/entity"
	"load-testing-proxy-server/internal/client/firebase_client"
	"log"
	"os"
	"sync"
	"time"
	"unsafe"
)

func FirebaseClientTest(concurrentUser int, receiveMessagePerClient int) {
	log.Printf("starting firebase load test with %d concurrent user and %d receive message per client", concurrentUser, receiveMessagePerClient)

	ctx := context.Background()

	var wg sync.WaitGroup

	avgTime := EndToEndResponseTime{
		sync: sync.Mutex{},
	}

	for i := 0; i < concurrentUser; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			go firebase_client.BusListener(ctx)
			go firebase_client.LocationListener(ctx)

			var totalRequest int
			for totalRequest < receiveMessagePerClient {
				msg := <-firebase_client.LocationBroadcaster

				if v, ok := firebase_client.ActiveBus[msg.BusID]; ok {
					resp := entity.BusLocationFirebaseResponse{
						Number:    v.Number,
						Plate:     v.Plate,
						Status:    v.Status,
						Route:     v.Route,
						IsActive:  v.IsActive,
						Heading:   msg.Heading,
						Latitude:  msg.Latitude,
						Longitude: msg.Longitude,
						Speed:     msg.Speed,
						Timestamp: msg.Timestamp,
					}

					responseTime := time.Since(resp.Timestamp)
					log.Printf("message received from firebase with latency: %s", responseTime)

					avgTime.sync.Lock()
					avgTime.times = append(avgTime.times, responseTime.Seconds())
					avgTime.sync.Unlock()

					totalRequest++
				}
			}
		}()
	}
	wg.Wait()

	log.Printf("median response time: %v second", median(avgTime.times))
}

func FirebaseDriverTest(concurrentUser int, sendMessagePerClient int) {
	log.Printf("starting firebase load test with %d concurrent user and %d send message per client", concurrentUser, sendMessagePerClient)

	throughput := Throughput{
		sync: sync.Mutex{},
	}

	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < concurrentUser; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			addr := os.Getenv("FIREBASE_BIKUNKU_SERVER")
			token := os.Getenv("FIREBASE_BIKUNKU_TOKEN")

			url := fmt.Sprintf("ws://%s/bus/streamfirebase?type=%s&token=%s", addr, "driver", token)
			dial, _, err := websocket_dialler.DefaultDialer.Dial(url, nil)
			if err != nil {
				log.Println(err)
				return
			}
			defer dial.Close()

			data := entity.SendLocationRequest{
				Long:    1,
				Lat:     1,
				Speed:   0,
				Heading: 0,
			}
			var totalRequest int
			for totalRequest < sendMessagePerClient {
				err = dial.WriteJSON(data)
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
