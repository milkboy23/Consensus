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
var done chan int

func main() {
	flag.Parse()
	hasToken = *hasTokenAtStart

	log.Printf("lPort %d | sPort %d | start %v", *listenPort, *sendPort, *hasTokenAtStart)

	go StartListener()
	StartSender()

	go RandomlyWantToken()
	<-done
}

func RandomlyWantToken() {
	for {
		log.Print("I want to use my token")

		if !hasToken {
			log.Print("BUT I DON'T HAVE IT >:(")
			for !hasToken {

			}
		}
		log.Print("I have it :D")

		log.Print("Using token!")
		time.Sleep(time.Second * 1)
		SendToken()

		sleepDuration := rand.Intn(10) + 1
		log.Printf("Sleep duration %v seconds", sleepDuration)
		time.Sleep(time.Duration(sleepDuration) * time.Second)
	}
}

func SendToken() {
	log.Print("Sending token!")
	_, err := sender.PassToken(context.Background(), &proto.Empty{})
	if err != nil {
		log.Fatalf("Error sending token | %v", err)
	}

	hasToken = false
	log.Print("Token has been sent.")
}

func (server *ConsensusPeerServer) PassToken(ctx context.Context, empty *proto.Empty) (*proto.Empty, error) {
	hasToken = true
	log.Print("We got the token!")

	return &proto.Empty{}, nil
}

func StartSender() {
	portString := fmt.Sprintf(":1600%d", *sendPort)
	dialOptions := grpc.WithTransportCredentials(insecure.NewCredentials())
	connection, connectionEstablishErr := grpc.NewClient(portString, dialOptions)
	if connectionEstablishErr != nil {
		log.Fatalf("Could not establish connection on port %s | %v", portString, connectionEstablishErr)
	}

	sender = proto.NewTokenRingClient(connection)
}

func StartListener() {
	portString := fmt.Sprintf(":1600%d", *listenPort)
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
