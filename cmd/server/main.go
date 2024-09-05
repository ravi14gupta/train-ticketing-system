package main

import (
	"log"
	"net"

	"github.com/ravi14gupta/train-ticketing-system/internal/ticket"
	pb "github.com/ravi14gupta/train-ticketing-system/proto"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":60051")
	if err != nil {
		log.Fatalf("train-ticketing-system failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTicketServiceServer(grpcServer, ticket.NewService())

	log.Printf("train-ticketing-system listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("train-ticketing-system failed to serve: %v", err)
	}
}
