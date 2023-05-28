package load_test

import (
	"context"
	"encoding/json"
	"fmt"
	"load-testing-proxy-server/entity"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
	"unsafe"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	websocket_dialler "github.com/fasthttp/websocket"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SafeActiveBus struct {
	V    map[int]entity.Bus
	Sync sync.Mutex
}

type SafeLocationBus struct {
	V    map[int]entity.BusLocationFirebaseResponse
	Sync sync.Mutex
}

var (
	LocationBroadcaster = make(chan entity.BusLocationFirebase)
	ActiveBus           = SafeActiveBus{
		V:    make(map[int]entity.Bus),
		Sync: sync.Mutex{},
	}
	StateBus = SafeLocationBus{
		V:    make(map[int]entity.BusLocationFirebaseResponse),
		Sync: sync.Mutex{},
	}
)

func Connection() *firebase.App {
	ctx := context.Background()
	sa := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Printf("error: %v", err)
	}
	return app
}

func FirebaseClientTestWithContext(ctx context.Context, concurrentUser int) {
	log.Printf("starting firebase load test with %d concurrent user", concurrentUser)

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

			// Connect Firebase
			app := Connection()

			// go firebase_client.BusListener(ctxWithTimeout, app)
			clientBus, err := app.Firestore(ctxWithTimeout)
			if err != nil {
				errorsOccur.HandleError(err)
				return
			}
			defer clientBus.Close()

			itBus := clientBus.Collection("buses").Snapshots(ctxWithTimeout)
			go func() {
				for {
					snap, err := itBus.Next()
					// DeadlineExceeded will be returned when ctx is cancelled.
					if status.Code(err) == codes.DeadlineExceeded {
						errorsOccur.HandleError(err)
						return
					}
					if err != nil {
						errorsOccur.HandleError(err)
						return
					}
					if snap != nil {
						for {
							doc, err := snap.Documents.Next()
							if err == iterator.Done {
								break
							}
							if err != nil {
								errorsOccur.HandleError(err)
								return
							}

							var data entity.Bus
							b, err := json.Marshal(doc.Data())
							if err != nil {
								errorsOccur.HandleError(err)
								return
							}

							err = json.Unmarshal(b, &data)
							if err != nil {
								errorsOccur.HandleError(err)
								return
							}

							data.ID, err = strconv.Atoi(doc.Ref.ID)
							if err != nil {
								errorsOccur.HandleError(err)
								return
							}

							// Adding/deleting to ActiveBus map
							if data.IsActive {
								ActiveBus.Sync.Lock()
								ActiveBus.V[data.ID] = data
								ActiveBus.Sync.Unlock()
							} else {
								ActiveBus.Sync.Lock()
								delete(ActiveBus.V, data.ID)
								ActiveBus.Sync.Unlock()
								StateBus.Sync.Lock()
								delete(StateBus.V, data.ID)
								StateBus.Sync.Unlock()
							}
						}
					}
				}
			}()

			time.Sleep(time.Second)

			// go firebase_client.LocationListener(ctxWithTimeout, app)
			clientLoc, err := app.Firestore(ctxWithTimeout)
			if err != nil {
				errorsOccur.HandleError(err)
				return
			}
			defer clientLoc.Close()

			itLoc := clientLoc.Collection("bus_locations").OrderBy("timestamp", firestore.Desc).Limit(1).Snapshots(ctxWithTimeout)
			for {
				snap, err := itLoc.Next()
				// DeadlineExceeded will be returned when ctx is cancelled.
				if status.Code(err) == codes.DeadlineExceeded {
					errorsOccur.HandleError(err)
					return
				}
				if err != nil {
					errorsOccur.HandleError(err)
					return
				}
				if snap != nil {
					for {
						doc, err := snap.Documents.Next()
						if err == iterator.Done {
							break
						}
						if err != nil {
							errorsOccur.HandleError(err)
							return
						}

						var msg entity.BusLocationFirebase
						b, err := json.Marshal(doc.Data())
						if err != nil {
							errorsOccur.HandleError(err)
							return
						}

						err = json.Unmarshal(b, &msg)
						if err != nil {
							errorsOccur.HandleError(err)
							return
						}

						// Update StateBus based on ActiveBus availability
						ActiveBus.Sync.Lock()
						v, ok := ActiveBus.V[msg.BusID]
						ActiveBus.Sync.Unlock()
						if ok {
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
							StateBus.Sync.Lock()
							StateBus.V[msg.BusID] = resp
							StateBus.Sync.Unlock()

							responseTime := time.Since(resp.Timestamp)

							avgTime.sync.Lock()
							avgTime.times = append(avgTime.times, responseTime.Seconds())
							avgTime.sync.Unlock()
						}
					}
				}
			}
		}()
	}

	<-ctx.Done()

	log.Printf("median response time: %v second", median(avgTime.times))
	log.Printf("error occured: %v", errorsOccur.errors)
}

func FirebaseClientTest(concurrentUser int, receiveMessagePerClient int) {
	log.Printf("starting firebase load test with %d concurrent user and %d receive message per client", concurrentUser, receiveMessagePerClient)

	ctx := context.Background()

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

			// Connect Firebase
			app := Connection()

			// go firebase_client.BusListener(ctx, app)
			clientBus, err := app.Firestore(ctx)
			if err != nil {
				errorsOccur.HandleError(err)
				return
			}
			defer clientBus.Close()

			itBus := clientBus.Collection("buses").Snapshots(ctx)
			go func() {
				for {
					snap, err := itBus.Next()
					// DeadlineExceeded will be returned when ctx is cancelled.
					if status.Code(err) == codes.DeadlineExceeded {
						errorsOccur.HandleError(err)
						return
					}
					if err != nil {
						errorsOccur.HandleError(err)
						return
					}
					if snap != nil {
						for {
							doc, err := snap.Documents.Next()
							if err == iterator.Done {
								break
							}
							if err != nil {
								errorsOccur.HandleError(err)
								return
							}

							var data entity.Bus
							b, err := json.Marshal(doc.Data())
							if err != nil {
								errorsOccur.HandleError(err)
								return
							}

							err = json.Unmarshal(b, &data)
							if err != nil {
								errorsOccur.HandleError(err)
								return
							}

							data.ID, err = strconv.Atoi(doc.Ref.ID)
							if err != nil {
								errorsOccur.HandleError(err)
								return
							}

							// Adding/deleting to ActiveBus map
							if data.IsActive {
								ActiveBus.Sync.Lock()
								ActiveBus.V[data.ID] = data
								ActiveBus.Sync.Unlock()
							} else {
								ActiveBus.Sync.Lock()
								delete(ActiveBus.V, data.ID)
								ActiveBus.Sync.Unlock()
								StateBus.Sync.Lock()
								delete(StateBus.V, data.ID)
								StateBus.Sync.Unlock()
							}
						}
					}
				}
			}()

			time.Sleep(time.Second)

			// go firebase_client.LocationListener(ctx, app)
			clientLoc, err := app.Firestore(ctx)
			if err != nil {
				errorsOccur.HandleError(err)
				return
			}
			defer clientLoc.Close()

			itLoc := clientLoc.Collection("bus_locations").OrderBy("timestamp", firestore.Desc).Limit(1).Snapshots(ctx)
			var totalRequest int
			for totalRequest < receiveMessagePerClient {
				snap, err := itLoc.Next()
				// DeadlineExceeded will be returned when ctx is cancelled.
				if status.Code(err) == codes.DeadlineExceeded {
					errorsOccur.HandleError(err)
					return
				}
				if err != nil {
					errorsOccur.HandleError(err)
					return
				}
				if snap != nil {
					for {
						doc, err := snap.Documents.Next()

						if err == iterator.Done {
							break
						}
						if err != nil {
							errorsOccur.HandleError(err)
							return
						}

						docData := doc.Data()

						var msg entity.BusLocationFirebase
						b, err := json.Marshal(docData)
						if err != nil {
							errorsOccur.HandleError(err)
							return
						}

						err = json.Unmarshal(b, &msg)
						if err != nil {
							errorsOccur.HandleError(err)
							return
						}

						// Update StateBus based on ActiveBus availability
						ActiveBus.Sync.Lock()
						v, ok := ActiveBus.V[msg.BusID]
						ActiveBus.Sync.Unlock()
						if ok {
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
							StateBus.Sync.Lock()
							StateBus.V[msg.BusID] = resp
							StateBus.Sync.Unlock()

							responseTime := time.Since(resp.Timestamp)

							avgTime.sync.Lock()
							avgTime.times = append(avgTime.times, responseTime.Seconds())
							avgTime.sync.Unlock()

							totalRequest++
						}
					}
				}
			}
		}()
	}
	wg.Wait()

	log.Printf("median response time: %v second", median(avgTime.times))
	log.Printf("error occured: %v", errorsOccur.errors)
}

// Deprecated: please use FirebaseDriverTest in integration package
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
