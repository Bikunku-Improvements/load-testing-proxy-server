package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	websocket_dialler "github.com/fasthttp/websocket"
)

func SendRouteFirebase(ctx context.Context, token string, route []byte) {
	addr := os.Getenv("FIREBASE_BIKUNKU_SERVER")

	url := fmt.Sprintf("ws://%s/bus/streamfirebase?type=%s&token=%s", addr, "driver", token)
	dial, _, err := websocket_dialler.DefaultDialer.Dial(url, nil)
	if err != nil {
		errorOccur.HandleError(err)
	}
	defer dial.Close()

	var routeReq []RouteData
	if err := json.Unmarshal(route, &routeReq); err != nil {
		errorOccur.HandleError(err)
	}

	for _, v := range routeReq {
		b, err := json.Marshal(v)
		if err != nil {
			errorOccur.HandleError(err)
			break
		}

		err = dial.WriteMessage(1, b)
		if err != nil {
			errorOccur.HandleError(err)
			break
		}

		err = throughput.AddSizeData(v)
		if err != nil {
			log.Printf("error when adding size data: %v", err)
			break
		}
		time.Sleep(time.Second)
	}
}

func LoginDriverFirebase(ctx context.Context, username, password string) (string, error) {
	addr := os.Getenv("FIREBASE_BIKUNKU_SERVER")

	url := fmt.Sprintf("http://%s/bus/loginAlt", addr)
	payload := map[string]interface{}{
		"username": username,
		"password": password,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("get %d response from server", resp.StatusCode)
	}

	var loginResp LoginResp
	if err = json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return "", err
	}

	return loginResp.Data.Token, nil
}

func FirebaseDriverTest(ctx context.Context) {
	now := time.Now()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		for _, v := range RedCredential {
			wg.Add(1)
			go func() {
				defer wg.Done()
				token, err := LoginDriverFirebase(ctx, v.Username, v.Password)
				if err != nil {
					log.Printf("failed login %s", v.Username)
					errorOccur.HandleError(err)
					return
				}
				log.Printf("success login %s", v.Username)
				SendRouteFirebase(ctx, token, Red)
			}()
			time.Sleep(1 * time.Second)
		}
	}()

	time.Sleep(1 * time.Second)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for _, v := range BlueCredential {
			wg.Add(1)
			go func() {
				defer wg.Done()
				token, err := LoginDriverFirebase(ctx, v.Username, v.Password)
				if err != nil {
					log.Printf("failed login %s", v.Username)
					errorOccur.HandleError(err)
					return
				}
				log.Printf("success login %s", v.Username)
				SendRouteFirebase(ctx, token, Blue)
			}()
			time.Sleep(1 * time.Second)
		}
	}()
	wg.Wait()

	log.Printf("total size data send: %v MB/s", float64(throughput.sizeData)/time.Since(now).Seconds()/1000000)
	log.Printf("total count data send: %v", throughput.totalData)
	log.Printf("error occured: %v", errorOccur.errors)
}
