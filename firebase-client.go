package main

import (
	"github.com/gofiber/websocket/v2"
	"log"
	"time"
)

func FirebaseClient(c *websocket.Conn) {
	for {
		time.Sleep(time.Second)
		if err := c.WriteMessage(1, []byte("hello-world")); err != nil {
			log.Println("write:", err)
			break
		}
	}
}
