package firebase_client

import (
	"github.com/gofiber/websocket/v2"
	"log"
	"time"
)

var clients = make(map[*websocket.Conn]bool)

func HandleConnection(c *websocket.Conn) {
	clients[c] = true
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			delete(clients, c)
			break
		}
	}
}

func HandleMessages() {
	for {
		// grab any next message from channel
		msg := <-LocationBroadcaster
		log.Printf("message received with latency: %s", time.Now().Sub(msg.Timestamp))
		if v, ok := ActiveBus[msg.BusID]; ok {
			messageClients(BusLocationResponse{
				Number:    v.Number,
				Plate:     v.Plate,
				Status:    v.Status,
				Route:     v.Route,
				IsActive:  v.IsActive,
				Heading:   msg.Heading,
				Latitude:  msg.Latitude,
				Longitude: msg.Longitude,
				Speed:     msg.Speed,
			})
		}
	}
}

func messageClients(msg BusLocationResponse) {
	// send to every client currently connected
	for client := range clients {
		messageClient(client, msg)
	}
}

func messageClient(client *websocket.Conn, msg BusLocationResponse) {
	err := client.WriteJSON(msg)
	if err != nil {
		log.Printf("error: %v", err)
		client.Close()
		delete(clients, client)
	}
}
