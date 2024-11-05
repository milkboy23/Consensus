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

var hasTokenAtStart = flag.Bool("start", false, "")
var listenPort = flag.Int("lPort", 0, "")
var sendPort = flag.Int("sPort", 1, "")

type ConsensusPeerServer struct {
	proto.UnimplementedTokenRingServer
}

var sender proto.TokenRingClient
var hasToken bool

func main() {
	flag.Parse()
	hasToken = *hasTokenAtStart

	StartListener()
	StartSender()

	go RandomlyWantToken()
}

func RandomlyWantToken() {
	for {
		fmt.Print("I want it :3")
		fmt.Print("BUT I DON'T HAVE IT >:(")
		for !hasToken {

		}

		if !hasToken {
			fmt.Print("This shit is broken.")
		}

		fmt.Print("Using token!")
		SendToken()

		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	}
}

func SendToken() {
	_, err := sender.PassToken(context.Background(), &proto.Empty{})
	if err != nil {
		log.Fatalf("Error sending token | %v", err)
	}

	hasToken = false
	fmt.Print("Token has been sent.")
}

func (server *ConsensusPeerServer) PassToken(ctx context.Context, empty *proto.Empty) (*proto.Empty, error) {
	hasToken = true

	fmt.Print("We got the token! Holding for 1 second...")
	time.Sleep(time.Second * 1)
	fmt.Print("Token passed.")

	return &proto.Empty{}, nil
}

func StartSender() {
	portString := fmt.Sprintf(":1600%d", sendPort)
	dialOptions := grpc.WithTransportCredentials(insecure.NewCredentials())
	connection, connectionEstablishErr := grpc.NewClient(portString, dialOptions)
	if connectionEstablishErr != nil {
		log.Fatalf("Could not establish connection on port %s | %v", portString, connectionEstablishErr)
	}

	sender = proto.NewTokenRingClient(connection)
}

func StartListener() {
	portString := fmt.Sprintf(":1600%d", listenPort)
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
