package ws_client

import (
	"encoding/json"
	"fmt"
	websocket_dialler "github.com/fasthttp/websocket"
	"github.com/gofiber/websocket/v2"
	"log"
	"os"
	"time"
)

func WSClient(c *websocket.Conn) {
	newestLocation := make(map[int]int)

	addr := os.Getenv("LEGACY_BIKUNKU_SERVER")
	dial, resp, err := websocket_dialler.DefaultDialer.Dial(fmt.Sprintf("ws://%s/bus/stream", addr), nil)
	if err != nil {
		log.Fatalln(resp, err)
	}
	defer dial.Close()
	defer c.Close()

	for {
		_, msg, err := dial.ReadMessage()
		if err != nil {
			log.Println("connection closed: ", err)
			break
		}

		var location []BusLocationResponse
		err = json.Unmarshal(msg, &location)
		if err != nil {
			log.Println("failed to unmarshall data: ", err)
			break
		}

		for _, v := range location {
			if _, ok := newestLocation[v.Id]; !ok {
				log.Printf("message received from websocket legacy with id=%d with latency: %s", v.Id, time.Now().Sub(v.Timestamp))
				newestLocation[v.Id] = v.Id
			}
		}

		err = c.WriteJSON(location)
		if err != nil {
			log.Println("failed to write: ", err)
			break
		}
	}

}
