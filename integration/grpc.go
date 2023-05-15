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

func SendRoute(ctx context.Context, token string, route []byte) {
	addr := os.Getenv("DRIVER_SERVICE_GRPC")

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		errorOccur.HandleError(err)
		return
	}

	md := metadata.New(map[string]string{
		"token": token,
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	svc := pb.NewLocationClient(conn)
	stream, err := svc.SendLocation(ctx)
	if err != nil {
		errorOccur.HandleError(err)
		return
	}

	var routeReq []pb.SendLocationRequest
	err = json.Unmarshal(route, &routeReq)
	if err != nil {
		errorOccur.HandleError(err)
		return
	}

	for _, v := range routeReq {
		err = stream.Send(&v)
		if err != nil {
			errorOccur.HandleError(err)
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
		wg.Add(1)
		for _, v := range RedCredential {
			go func() {
				defer wg.Done()
				token, err := LoginDriver(ctx, v.Username, v.Password)
				if err != nil {
					log.Printf("failed login %s", v.Username)
					errorOccur.HandleError(err)
					return
				}
				log.Printf("success login %s", v.Username)
				SendRoute(ctx, token, Red)
			}()
			time.Sleep(5 * time.Second)
		}
		wg.Done()
	}()

	time.Sleep(2 * time.Second)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, v := range BlueCredential {
			go func() {
				wg.Add(1)
				defer wg.Done()
				token, err := LoginDriver(ctx, v.Username, v.Password)
				if err != nil {
					log.Printf("failed login %s", v.Username)
					errorOccur.HandleError(err)
					return
				}
				log.Printf("success login %s", v.Username)
				SendRoute(ctx, token, Blue)
			}()
			time.Sleep(5 * time.Second)
		}
	}()
	wg.Wait()
}
