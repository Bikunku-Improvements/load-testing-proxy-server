package load_test

import (
	"context"
	"load-testing-proxy-server/entity"
	"load-testing-proxy-server/internal/client/firebase_client"
	"log"
	"sync"
	"time"
)

func FirebaseTest(concurrentUser int, receiveMessagePerClient int) {
	log.Printf("starting firebase load test with %d concurrent user and %d receive message per client", concurrentUser, receiveMessagePerClient)

	ctx := context.Background()

	var wg sync.WaitGroup

	avgTime := AverageTime{
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
					avgTime.count++
					avgTime.sync.Unlock()

					totalRequest++
				}
			}
		}()
	}
	wg.Wait()

	log.Printf("median response time: %v second", median(avgTime.times))
}
