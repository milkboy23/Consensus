package main

import (
	proto "Consensus/GRPC"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"math/rand"
	"net"
	"time"
)

var starter = flag.Bool("s", false, "")
var id = flag.Int("id", 0, "")

type ConsensusPeerServer struct {
	proto.UnimplementedTokenRingServer
}

var sender proto.TokenRingClient
var wantToken bool

func main() {
	flag.Parse() // Parses program arguments into flag variables

	StartSender()
	StartListener()

	if *starter { // If this node is the starter, inject a token into the system by passing without explicitly having received one
		go PassToken()
	}
	RandomlyWantToken() // Randomly set wantToken to true every once in a while
}

func RandomlyWantToken() {
	for {
		sleepDuration := rand.Intn(12) + 3                     // Between 3 and 15 seconds
		time.Sleep(time.Duration(sleepDuration) * time.Second) // Wait until we want the token again

		wantToken = true
		log.Print("I want to use my token!") // If you don't have a print statement here, the whole program breaks :)

		for wantToken {
		} // Wait until we no longer want token (after getting and using token)
	}
}

func (server *ConsensusPeerServer) ReceiveToken(ctx context.Context, empty *proto.Empty) (*proto.Empty, error) {
	// We just received a token (access to CS),
	// and do different things depending on whether we want to use it right now
	if !wantToken {
		log.Print("Just got the token, passing...")
		time.Sleep(time.Second) // Hold the token for a second
		go PassToken()          // Pass token to next node because we're done with it
	} else {
		log.Print("Just got the token, using...")
		UseToken() // Hold the token for 3 seconds
	}
	return &proto.Empty{}, nil
}

func UseToken() {
	time.Sleep(time.Second * 3) // Simulate work/accessing CS

	log.Print("Done using token.")
	wantToken = false // We no longer want the token

	go PassToken() // Pass token to next node because we're done with it
}

func PassToken() {
	_, err := sender.ReceiveToken(context.Background(), &proto.Empty{}) // Send the token via. gRPC
	// P.S. This ^^^ gRPC call is blocking, because it waits for a response from the receiving node.
	// Which is why we throw it to a goroutine.
	if err != nil {
		log.Fatalf("Error sending token | %v", err)
	}
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

	go serveListener(grpcListener, listener)
	// Throw listener serving method to goroutine,
	// because otherwise the program just waits/stops here, and that would be annoying
}

func serveListener(grpcListener *grpc.Server, listener net.Listener) {
	serveListenerErr := grpcListener.Serve(listener) // Activate listener (blocking call, program waits here until shutdown)
	if serveListenerErr != nil {
		log.Fatalf("Failed to serve listener | %v", serveListenerErr)
	}
}
