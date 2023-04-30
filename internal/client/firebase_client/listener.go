package firebase_client

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"strconv"
)

func BusListener(ctx context.Context) {
	sa := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Printf("error: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Printf("error: %v", err)
	}
	defer client.Close()

	it := client.Collection("buses").Snapshots(ctx)
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

				var data Bus
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
					ActiveBus[data.ID] = data
				} else {
					delete(ActiveBus, data.ID)
				}
			}
		}
	}
}

func LocationListener(ctx context.Context) {
	sa := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Printf("error: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Printf("error: %v", err)
	}
	defer client.Close()

	it := client.Collection("bus_locations").OrderBy("timestamp", firestore.Desc).Limit(1).Snapshots(ctx)
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
				
				var msg BusLocation
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
