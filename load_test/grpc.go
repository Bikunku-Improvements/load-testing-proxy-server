package load_test

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"load-testing-proxy-server/entity"
	"load-testing-proxy-server/grpc/pb"
	"log"
	"os"
	"sync"
	"time"
)

func GRPCTest(concurrentUser int, receiveMessagePerClient int) {
	log.Printf("starting grpc load test with %d concurrent user and %d receive message per client", concurrentUser, receiveMessagePerClient)

	var wg sync.WaitGroup

	avgTime := AverageTime{
		sync: sync.Mutex{},
	}

	log.Printf("adding client...")
	for i := 0; i < concurrentUser; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			addr := os.Getenv("PASSENGER_SERVICE_ADDR")

			conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			defer conn.Close()
			if err != nil {
				log.Printf("error connecting to grpc: %v", err)
				return
			}

			client := pb.NewLocationClient(conn)
			stream, err := client.SubscribeLocation(context.Background(), &pb.SubscribeLocationRequest{})
			if err != nil {
				log.Printf("error connecting to grpc: %v", err)
				return
			}

			// receive message
			var totalRequest int
			for totalRequest < receiveMessagePerClient {
				location, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Printf("%v.SubscribeLocation(_) = _, %v", client, err)
					return
				}

				b, err := json.Marshal(location)
				if err != nil {
					log.Printf("failed to unmarshall: %v", err)
				}

				var loc entity.BusLocationGRPC
				err = json.Unmarshal(b, &loc)
				if err != nil {
					log.Printf("failed to unmarshall: %v", err)
				}

				responseTime := time.Since(loc.CreatedAt)
				log.Printf("message received from grpc with latency: %s", responseTime)

				avgTime.sync.Lock()
				avgTime.times = append(avgTime.times, responseTime.Seconds())
				avgTime.count++
				avgTime.sync.Unlock()

				totalRequest++
			}

		}()
	}
	wg.Wait()

	log.Printf("median response time: %v second", median(avgTime.times))
}
