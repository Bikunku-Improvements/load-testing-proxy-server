package main

import (
	"github.com/joho/godotenv"
	"load-testing-proxy-server/load_test"
	"log"
	"os"
	"strconv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("error loading .env file: %v", err)
	}

	args := os.Args[1:]
	if len(args) == 1 && args[0] == "help" {
		log.Println("Usage go run main.go [type] [concurrent user] [receive message per client]")
		log.Println("type options grpc, ws-legacy, firebase")
	}

	if len(args) < 3 {
		log.Fatalf("arguments need to be 'type concurentUser receiveMessagePerClient'")
	}

	testType := args[0]
	concurrentUser, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatalf("concurrent user needs to be integer")
	}

	receiveMessagePerClient, err := strconv.Atoi(args[2])
	if err != nil {
		log.Fatalf("receive message per client needs to be integer")
	}

	switch testType {
	case "grpc":
		load_test.GRPCTest(concurrentUser, receiveMessagePerClient)
	case "firebase":
		load_test.FirebaseTest(concurrentUser, receiveMessagePerClient)
	case "ws-legacy":
		load_test.WSTest(concurrentUser, receiveMessagePerClient)
	default:
		log.Fatalf("type need to be grpc, firebase, or ws-legacy")
	}
}
