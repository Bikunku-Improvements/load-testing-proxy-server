package grpc_client

import (
	"context"
	"github.com/gofiber/websocket/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"load-testing-proxy-server/grpc/pb"
	"log"
	"os"
	"time"
)

func GRPCClient(c *websocket.Conn) {
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
		return
	}

	for {
		location, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.SubscribeLocation(_) = _, %v", client, err)
		}

		parsedTime, err := time.Parse(time.RFC3339, location.GetCreatedAt())
		if err != nil {
			log.Fatalf("%v.SubscribeLocation(_) = _, %v", client, err)
		}
		log.Printf("message received from grpc with latency: %s", time.Now().Sub(parsedTime))
		c.WriteJSON(location)
	}
}
