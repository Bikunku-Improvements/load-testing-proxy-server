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

func GRPCClientTestWithContext(ctx context.Context, concurrentUser int) {
	log.Printf("starting grpc load test with %d concurrent user", concurrentUser)

	avgTime := EndToEndResponseTime{
		sync: sync.Mutex{},
	}

	errorsOccur := ErrorOccur{
		errors: make(map[string]int),
		sync:   sync.Mutex{},
	}

	for i := 0; i < concurrentUser; i++ {
		go func() {
			ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Minute)
			defer cancel()
			addr := os.Getenv("PASSENGER_SERVICE_ADDR")

			conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			defer conn.Close()
			if err != nil {
				errorsOccur.HandleError(err)
				return
			}

			client := pb.NewLocationClient(conn)
			stream, err := client.SubscribeLocation(ctxWithTimeout, &pb.SubscribeLocationRequest{})
			if err != nil {
				errorsOccur.HandleError(err)
				return
			}

			// receive message
			for {
				location, err := stream.Recv()
				if err == io.EOF {
					break
				}

				if err != nil {
					errorsOccur.HandleError(err)
					return
				}

				b, err := json.Marshal(location)
				if err != nil {
					errorsOccur.HandleError(err)
					return
				}

				var loc entity.BusLocationGRPC
				err = json.Unmarshal(b, &loc)
				if err != nil {
					errorsOccur.HandleError(err)
					return
				}

				responseTime := time.Since(loc.CreatedAt)
				log.Printf("message received from grpc with latency: %s", responseTime)

				avgTime.sync.Lock()
				avgTime.times = append(avgTime.times, responseTime.Seconds())
				avgTime.sync.Unlock()
			}

		}()
	}

	<-ctx.Done()

	log.Printf("median response time: %v second", median(avgTime.times))
	log.Printf("error occured: %v", errorsOccur.errors)
}

func GRPCClientTest(concurrentUser, receiveMessagePerClient int) {
	log.Printf("starting grpc load test with %d concurrent user and %d receive message per client", concurrentUser, receiveMessagePerClient)
	var wg sync.WaitGroup

	avgTime := EndToEndResponseTime{
		sync: sync.Mutex{},
	}

	errorsOccur := ErrorOccur{
		errors: make(map[string]int),
		sync:   sync.Mutex{},
	}

	for i := 0; i < concurrentUser; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			addr := os.Getenv("PASSENGER_SERVICE_ADDR")

			conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			defer conn.Close()
			if err != nil {
				errorsOccur.HandleError(err)
				return
			}

			client := pb.NewLocationClient(conn)
			stream, err := client.SubscribeLocation(context.Background(), &pb.SubscribeLocationRequest{})
			if err != nil {
				errorsOccur.HandleError(err)
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
					errorsOccur.HandleError(err)
					return
				}

				b, err := json.Marshal(location)
				if err != nil {
					errorsOccur.HandleError(err)
					return
				}

				var loc entity.BusLocationGRPC
				err = json.Unmarshal(b, &loc)
				if err != nil {
					errorsOccur.HandleError(err)
					return
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
	log.Printf("error occured: %v", errorsOccur.errors)
}

// Deprecated: please use GRPCDriverTest in integration package
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

			svc := pb.NewLocationClient(conn)
			stream, err := svc.SendLocation(ctx)
			if err != nil {
				log.Printf("error connecting to grpc: %v", err)
				return
			}
			var totalRequest int
			for totalRequest < sendMessagePerClient {
				data := []byte(`{"long": 91, "lat": 20, "speed": 0.00, "heading": 0}`)

				var dataMarshalled pb.SendLocationRequest
				err = json.Unmarshal(data, &dataMarshalled)
				if err != nil {
					log.Printf("error when marhsal input data: %v", err)
					return
				}

				err = stream.Send(&dataMarshalled)
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
