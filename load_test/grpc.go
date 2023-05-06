package load_test

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"io"
	"load-testing-proxy-server/entity"
	"load-testing-proxy-server/grpc/pb"
	"log"
	"os"
	"sync"
	"time"
	"unsafe"
)

func GRPCClientTest(concurrentUser, receiveMessagePerClient int) {
	log.Printf("starting grpc load test with %d concurrent user and %d receive message per client", concurrentUser, receiveMessagePerClient)

	avgTime := EndToEndResponseTime{
		sync: sync.Mutex{},
	}

	var wg sync.WaitGroup
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
				avgTime.sync.Unlock()

				totalRequest++
			}

		}()
	}
	wg.Wait()

	log.Printf("median response time: %v second", median(avgTime.times))
}

func GRPCDriverTest(concurrentUser, sendMessagePerClient int) {
	log.Printf("starting grpc driver load test with %d concurrent user and %d send message per client", concurrentUser, sendMessagePerClient)

	throughput := Throughput{
		sync: sync.Mutex{},
	}

	ctx := context.Background()

	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < concurrentUser; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			addr := os.Getenv("DRIVER_SERVICE_GRPC")

			conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			defer conn.Close()
			if err != nil {
				log.Printf("error connecting to grpc: %v", err)
				return
			}

			md := metadata.New(map[string]string{
				"token": os.Getenv("DRIVER_SERVICE_TOKEN"),
			})
			ctx = metadata.NewOutgoingContext(ctx, md)

			data := pb.SendLocationRequest{
				Long:    1,
				Lat:     1,
				Speed:   0,
				Heading: 0,
			}

			svc := pb.NewLocationClient(conn)
			stream, err := svc.SendLocation(ctx)
			if err != nil {
				log.Printf("error connecting to grpc: %v", err)
				return
			}
			var totalRequest int
			for totalRequest < sendMessagePerClient {
				err = stream.Send(&data)
				if err != nil {
					log.Printf("error sending data to grpc: %v", err)
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
