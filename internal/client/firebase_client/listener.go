package firebase_client

import (
	"context"
	"encoding/json"
	"fmt"
	"load-testing-proxy-server/entity"
	"log"
	"os"
	"strconv"
	"sync"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SafeActiveBus struct {
	V map[int]entity.Bus
	Sync sync.Mutex
}

var (
	LocationBroadcaster = make(chan entity.BusLocationFirebase)
	ActiveBus						= SafeActiveBus{
		V: make(map[int]entity.Bus),
		Sync: sync.Mutex{},
	}
)

func Connection() (*firebase.App) {
	ctx := context.Background()
	sa := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Printf("error: %v", err)
	}
	return app
}

func BusListener(ctx context.Context, app *firebase.App) {
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Printf("error: %v", err)
	}
	defer client.Close()

	it := client.Collection("buses").Snapshots(ctx)
	defer it.Stop()
	for {
		snap, err := it.Next()
		// DeadlineExceeded will be returned when ctx is cancelled.
		if status.Code(err) == codes.DeadlineExceeded {
			log.Printf("failed code deadline exceeded: %v", err)
		}
		if err != nil {
			log.Printf("Snapshots.Next: %v", err)
		}
		if snap != nil {
			for {
				doc, err := snap.Documents.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					fmt.Printf("Documents.Next: %v", err)
				}

				var data entity.Bus
				b, err := json.Marshal(doc.Data())
				if err != nil {
					log.Printf("error: %v", err)
				}

				err = json.Unmarshal(b, &data)
				if err != nil {
					log.Printf("error: %v", err)
				}

				data.ID, err = strconv.Atoi(doc.Ref.ID)
				if err != nil {
					log.Printf("invalid parsing id: %v", err)
				}

				if data.IsActive {
					ActiveBus.Sync.Lock()
					ActiveBus.V[data.ID] = data
					ActiveBus.Sync.Unlock()
				} else {
					ActiveBus.Sync.Lock()
					delete(ActiveBus.V, data.ID)
					ActiveBus.Sync.Unlock()
				}
			}
		}
	}
}

func LocationListener(ctx context.Context, app *firebase.App) {
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Printf("error: %v", err)
	}
	defer client.Close()

	it := client.Collection("bus_locations").OrderBy("timestamp", firestore.Desc).Limit(1).Snapshots(ctx)
	defer it.Stop()
	for {
		snap, err := it.Next()
		// DeadlineExceeded will be returned when ctx is cancelled.
		if status.Code(err) == codes.DeadlineExceeded {
			log.Printf("failed code deadline exceeded: %v", err)
		}
		if err != nil {
			log.Printf("Snapshots.Next: %v", err)
		}
		if snap != nil {
			for {
				doc, err := snap.Documents.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					fmt.Printf("Documents.Next: %v", err)
				}

				var msg entity.BusLocationFirebase
				b, err := json.Marshal(doc.Data())
				if err != nil {
					log.Printf("error: %v", err)
				}

				err = json.Unmarshal(b, &msg)
				if err != nil {
					log.Printf("error: %v", err)
				}

				LocationBroadcaster <- msg
			}
		}
	}
}
