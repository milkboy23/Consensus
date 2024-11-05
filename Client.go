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
var hasToken bool

func main() {
	flag.Parse()
	hasToken = *starter

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
		time.Sleep(time.Second * 1) // Simulate work/accessing CS
		log.Print("Done using token.")

		PassToken()
	}
}

func PassToken() {
	log.Print("Passing token!")
	_, err := sender.ReceiveToken(context.Background(), &proto.Empty{})
	if err != nil {
		log.Fatalf("Error sending token | %v", err)
	}

	hasToken = false
	log.Print("Token has been passed.")
}

func (server *ConsensusPeerServer) ReceiveToken(ctx context.Context, empty *proto.Empty) (*proto.Empty, error) {
	if hasToken {
		log.Fatalf("Bro... I already have a token, wtf are you doing :|")
	}

	hasToken = true
	log.Print("We just got the token!")

	return &proto.Empty{}, nil
}

func StartSender() {
	port := *id + 1
	if *starter {
		port = 0
	}

	portString := fmt.Sprintf(":1600%d", port)
	dialOptions := grpc.WithTransportCredentials(insecure.NewCredentials())
	connection, connectionEstablishErr := grpc.NewClient(portString, dialOptions)
	if connectionEstablishErr != nil {
		log.Fatalf("Could not establish connection on port %s | %v", portString, connectionEstablishErr)
	}

	sender = proto.NewTokenRingClient(connection)
}

func StartListener() {
	port := *id
	portString := fmt.Sprintf(":1600%d", port)
	listener, listenErr := net.Listen("tcp", portString)
	if listenErr != nil {
		log.Fatalf("Failed to listen on port %s | %v", portString, listenErr)
	}

	grpcListener := grpc.NewServer()
	proto.RegisterTokenRingServer(grpcListener, &ConsensusPeerServer{})

	serveListenerErr := grpcListener.Serve(listener)
	if serveListenerErr != nil {
		log.Fatalf("Failed to serve listener | %v", serveListenerErr)
	}
}
