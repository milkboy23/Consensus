package main

import (
	proto "Consensus/GRPC"
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var starter = flag.Bool("s", false, "")
var starterWants = flag.Bool("want", false, "")
var id = flag.Int("id", 0, "")

type ConsensusPeerServer struct {
	proto.UnimplementedTokenRingServer
}

var sender proto.TokenRingClient
var hasToken bool
var wantsToken bool

func main() {
	flag.Parse()        // Parses program arguments into flag variables
	hasToken = *starter // Set from flags if this node has a token on start
	wantsToken = *starterWants

	go RandomlyWantToken()

	StartSender()
	StartListener()
}

func RandomlyWantToken() {

	val := rand.Intn(1)
	if val == 0 {
		wantsToken = true
	} else if val == 1 {
		wantsToken = false
	}

	if wantsToken {
		for {
			sleepDuration := rand.Intn(10) + 1
			time.Sleep(time.Duration(sleepDuration) * time.Second) // Wait until we want the token
			log.Print("I want to use my token")

			if !hasToken {
				log.Print("BUT I DON'T HAVE IT >:(")
				for !hasToken {
				} // Wait until we get the token
			}

			log.Print("I have it :D, using token!")
			time.Sleep(time.Second * 1) // Simulate work/accessing CS after receiving
			log.Print("Done using token.")

			PassToken() // Pass token to next node
		}
	} else {
		PassToken()
	}
}

func PassToken() {
	log.Print("Passing token!")
	_, err := sender.ReceiveToken(context.Background(), &proto.Empty{}) // Send the token via. gRPC
	if err != nil {
		log.Fatalf("Error sending token | %v", err)
	}

	hasToken = false // This node no longer has the token
	log.Print("Token has been passed.")
}

func (server *ConsensusPeerServer) ReceiveToken(ctx context.Context, empty *proto.Empty) (*proto.Empty, error) {
	if hasToken {
		log.Fatalf("Bro... I already have a token, wtf are you doing :|")
	}

	hasToken = true // Node received the token through gRPC
	log.Print("We just got the token!")

	return &proto.Empty{}, nil
}

func StartSender() {
	port := *id + 1 // Neighbour's port
	if *starter {   // If it's the starter, neighbour must be first node and have id 0
		port = 0
	}

	portString := fmt.Sprintf(":1600%d", port) // Format port string
	dialOptions := grpc.WithTransportCredentials(insecure.NewCredentials())
	connection, connectionEstablishErr := grpc.NewClient(portString, dialOptions) // Register connection to neighbour node, so we can send commands to it
	if connectionEstablishErr != nil {
		log.Fatalf("Could not establish connection on port %s | %v", portString, connectionEstablishErr)
	}

	sender = proto.NewTokenRingClient(connection) // Sender variable is set to be a token ring node with the connection to its port
}

func StartListener() {
	port := *id                                          // Listen on our own port
	portString := fmt.Sprintf(":1600%d", port)           // Format port string
	listener, listenErr := net.Listen("tcp", portString) // Listen with TCP on the 1600[id] port
	if listenErr != nil {
		log.Fatalf("Failed to listen on port %s | %v", portString, listenErr)
	}

	grpcListener := grpc.NewServer() // Get reference to gRPC server (listener from which commands/gRPC calls come through)
	proto.RegisterTokenRingServer(grpcListener, &ConsensusPeerServer{})

	serveListenerErr := grpcListener.Serve(listener) // Activate listener (blocking call, program waits here until shutdown)
	if serveListenerErr != nil {
		log.Fatalf("Failed to serve listener | %v", serveListenerErr)
	}
}
