package main

import (
	"github.com/joho/godotenv"
	"load-testing-proxy-server/integration"
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
		log.Println("Usage go run main.go [type] [endpoint_type] [concurrent user] [receive/send message per client]")
		log.Println("type options are driver or client")
		log.Println("endpoint type options are grpc, ws-legacy, firebase")
	}
	//
	//if len(args) < 4 {
	//	log.Fatalf("arguments need to be '[type] [endpoint_type] [concurrent user] [receive/send message per client]'")
	//}

	clientType := args[0]
	switch clientType {
	case "driver":
		endpointType := args[1]
		concurrentUser, err := strconv.Atoi(args[2])
		if err != nil {
			log.Fatalf("concurrent user needs to be integer")
		}

		sendMessagePerClient, err := strconv.Atoi(args[3])
		if err != nil {
			log.Fatalf("send message per client needs to be integer")
		}

		switch endpointType {
		case "grpc":
			load_test.GRPCDriverTest(concurrentUser, sendMessagePerClient)
		case "firebase":
			load_test.FirebaseDriverTest(concurrentUser, sendMessagePerClient)
		case "ws-legacy":
			load_test.WSLegacyDriverTest(concurrentUser, sendMessagePerClient)
		default:
			log.Fatalf("type should be grpc, firebase, or ws-legacy")
		}
	case "client":
		endpointType := args[1]
		concurrentUser, err := strconv.Atoi(args[2])
		if err != nil {
			log.Fatalf("concurrent user needs to be integer")
		}

		receiveMessagePerClient, err := strconv.Atoi(args[3])
		if err != nil {
			log.Fatalf("receive message per client needs to be integer")
		}

		switch endpointType {
		case "grpc":
			load_test.GRPCClientTest(concurrentUser, receiveMessagePerClient)
		case "firebase":
			load_test.FirebaseClientTest(concurrentUser, receiveMessagePerClient)
		case "ws-legacy":
			load_test.WSLegacyClientTest(concurrentUser, receiveMessagePerClient)
		default:
			log.Fatalf("type need to be grpc, firebase, or ws-legacy")
		}
	case "integration":
		switch args[1] {
		case "grpc":
			integration.GRPCDriverTest()
		case "ws-legacy":
			integration.WSLegacyDriverTest()
		}
	default:
		log.Fatalf("type must be client or driver")
	}
}
