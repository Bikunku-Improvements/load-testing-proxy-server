package integration

import (
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"load-testing-proxy-server/grpc/pb"
	"log"
	"os"
	"sync"
	"time"
)

func SendRoute(ctx context.Context, wg *sync.WaitGroup, token string, route []byte) {
	defer wg.Done()
	addr := os.Getenv("DRIVER_SERVICE_GRPC")

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		log.Printf("error connecting to grpc: %v", err)
		return
	}

	md := metadata.New(map[string]string{
		"token": token,
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	svc := pb.NewLocationClient(conn)
	stream, err := svc.SendLocation(ctx)
	if err != nil {
		log.Printf("error connecting to grpc: %v", err)
		return
	}

	var routeReq []pb.SendLocationRequest
	err = json.Unmarshal(route, &routeReq)
	if err != nil {
		log.Printf("error when marhsal input data: %v", err)
		return
	}

	for _, v := range routeReq {
		err = stream.Send(&v)
		if err != nil {
			log.Printf("error sending data to grpc: %v", err)
			return
		}

		time.Sleep(time.Second)
	}
}

func LoginDriver(ctx context.Context, username, password string) (string, error) {
	addr := os.Getenv("DRIVER_SERVICE_GRPC")

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		log.Printf("error connecting to grpc: %v", err)
		return "", err
	}

	svc := pb.NewUserClient(conn)
	login, err := svc.Login(ctx, &pb.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return "", err
	}

	return login.Token, nil
}

func GRPCDriverTest() {
	ctx := context.Background()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		for _, v := range RedCredential {
			go func() {
				token, err := LoginDriver(ctx, v.Username, v.Password)
				if err != nil {
					log.Printf("unable to log in: %v", err)
					return
				}

				SendRoute(ctx, &wg, token, Red)
			}()
			time.Sleep(5 * time.Second)
		}
	}()

	time.Sleep(2 * time.Second)

	wg.Add(1)
	go func() {
		for _, v := range BlueCredential {
			go func() {
				token, err := LoginDriver(ctx, v.Username, v.Password)
				if err != nil {
					log.Printf("unable to log in: %v", err)
					return
				}
				SendRoute(ctx, &wg, token, Blue)
			}()
			time.Sleep(5 * time.Second)
		}
	}()
	wg.Wait()
}
