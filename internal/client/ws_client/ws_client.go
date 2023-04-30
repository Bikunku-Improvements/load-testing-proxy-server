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
	addr := os.Getenv("LEGACY_BIKUNKU_SERVER")
	dial, resp, err := websocket_dialler.DefaultDialer.Dial(fmt.Sprintf("wss://%s/bus/stream", addr), nil)
	if err != nil {
		log.Fatalln(resp, err)
	}
	defer dial.Close()

	for {
		_, msg, err := dial.ReadMessage()
		if err != nil {
			log.Println("connection closed: ", err)
			return
		}

		var location []BusLocationResponse
		err = json.Unmarshal(msg, &location)
		if err != nil {
			log.Println("failed to unmarshall data: ", err)
			return
		}

		for _, v := range location {
			if v.IsNewLocation {
				log.Printf("message received from websocket legacy with latency: %s", time.Now().Sub(v.Timestamp))
			}
		}

		err = c.WriteJSON(location)
		if err != nil {
			log.Println("failed to write: ", err)
		}
	}

}
