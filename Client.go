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
var id = flag.Int("id", 0, "")

type ConsensusPeerServer struct {
	proto.UnimplementedTokenRingServer
}

var sender proto.TokenRingClient
var hasToken bool

func main() {
	flag.Parse()        // TODO
	hasToken = *starter // set to the starter's bool

	go StartListener()
	StartSender()

	RandomlyWantToken()
}

func RandomlyWantToken() {
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

		PassToken()
	}
}

func PassToken() {
	log.Print("Passing token!")
	_, err := sender.ReceiveToken(context.Background(), &proto.Empty{}) // send the token
	if err != nil {
		log.Fatalf("Error sending token | %v", err)
	}

	hasToken = false // this node no longer has the token
	log.Print("Token has been passed.")
}

func (server *ConsensusPeerServer) ReceiveToken(ctx context.Context, empty *proto.Empty) (*proto.Empty, error) {
	if hasToken {
		log.Fatalf("Bro... I already have a token, wtf are you doing :|")
	}

	hasToken = true // node got the token
	log.Print("We just got the token!")

	return &proto.Empty{}, nil
}

func StartSender() {
	port := *id + 1 // neighbour's port
	if *starter {   // if it's the starter, neighbour's id is 0
		port = 0
	}

	portString := fmt.Sprintf(":1600%d", port)                                    // format
	dialOptions := grpc.WithTransportCredentials(insecure.NewCredentials())       // ... TODO
	connection, connectionEstablishErr := grpc.NewClient(portString, dialOptions) // register new node to gRPC server and get connection
	if connectionEstablishErr != nil {
		log.Fatalf("Could not establish connection on port %s | %v", portString, connectionEstablishErr)
	}

	sender = proto.NewTokenRingClient(connection) // sender is set to be a token ring node with the connection to its port
}

func StartListener() {
	port := *id                                          // listen on own port
	portString := fmt.Sprintf(":1600%d", port)           // format to always accessible port
	listener, listenErr := net.Listen("tcp", portString) // listener to listen with tcp on the 1600id port
	if listenErr != nil {
		log.Fatalf("Failed to listen on port %s | %v", portString, listenErr)
	}

	grpcListener := grpc.NewServer() // make gRPC server (en ny gRPC server hver gang??)
	proto.RegisterTokenRingServer(grpcListener, &ConsensusPeerServer{})

	serveListenerErr := grpcListener.Serve(listener) // give the listener to the gRPC server
	if serveListenerErr != nil {
		log.Fatalf("Failed to serve listener | %v", serveListenerErr)
	}
}
