package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"
	"load-testing-proxy-server/internal/client/firebase_client"
	"load-testing-proxy-server/internal/client/grpc_client"
	"load-testing-proxy-server/internal/client/ws_client"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("error loading .env file: %v", err)
	}

	ctx := context.Background()
	app := fiber.New()

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/firebase-client", websocket.New(firebase_client.HandleConnection))
	app.Get("/grpc-client", websocket.New(grpc_client.GRPCClient))
	app.Get("/ws-client", websocket.New(ws_client.WSClient))

	go firebase_client.HandleMessages()
	go firebase_client.LocationListener(ctx)
	go firebase_client.BusListener(ctx)

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}
