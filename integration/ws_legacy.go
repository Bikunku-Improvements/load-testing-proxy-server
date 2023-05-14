package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	websocket_dialler "github.com/fasthttp/websocket"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type LoginResp struct {
	Status string `json:"status"`
	Error  string `json:"error"`
	Data   struct {
		Id       int    `json:"id"`
		Number   int    `json:"number"`
		Plate    string `json:"plate"`
		Status   string `json:"status"`
		Route    string `json:"route"`
		IsActive bool   `json:"isActive"`
		Token    string `json:"token"`
	} `json:"data"`
}

func SendRouteWSLegacy(ctx context.Context, wg *sync.WaitGroup, token string, route []byte) {
	defer wg.Done()
	addr := os.Getenv("LEGACY_BIKUNKU_SERVER")

	url := fmt.Sprintf("ws://%s/bus/stream?type=%s&token=%s", addr, "driver", token)
	dial, resp, err := websocket_dialler.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Println(resp, err)
		return
	}
	defer dial.Close()

	var routeReq []RouteData
	if err := json.Unmarshal(route, &routeReq); err != nil {
		log.Printf("error when unmarshal data: %v\n", err)
		return
	}

	for _, v := range routeReq {
		b, err := json.Marshal(v)
		if err != nil {
			log.Printf("error when marshal: %v", err)
			return
		}

		err = dial.WriteMessage(1, b)
		if err != nil {
			log.Println(err)
			return
		}

		time.Sleep(time.Second)
	}
}

func LoginDriverWSLegacy(ctx context.Context, username, password string) (string, error) {
	addr := os.Getenv("LEGACY_BIKUNKU_SERVER")

	url := fmt.Sprintf("http://%s/bus/login", addr)
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

func WSLegacyDriverTest() {
	ctx := context.Background()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		for _, v := range RedCredential {
			go func() {
				token, err := LoginDriverWSLegacy(ctx, v.Username, v.Password)
				if err != nil {
					log.Printf("unable to log in: %v", err)
					return
				}
				SendRouteWSLegacy(ctx, &wg, token, Red)
			}()
			time.Sleep(5 * time.Second)
		}
	}()

	time.Sleep(2 * time.Second)

	wg.Add(1)
	go func() {
		for _, v := range BlueCredential {
			go func() {
				token, err := LoginDriverWSLegacy(ctx, v.Username, v.Password)
				if err != nil {
					log.Printf("unable to log in: %v", err)
					return
				}
				SendRouteWSLegacy(ctx, &wg, token, Blue)
			}()
			time.Sleep(5 * time.Second)
		}
	}()
	wg.Wait()
}
