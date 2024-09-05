package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/ravi14gupta/train-ticketing-system/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	// Establish the gRPC connection
	conn, err := grpc.NewClient("localhost:60051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	defer conn.Close()
	client := pb.NewTicketServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 1. Purchase a ticket for Ravi Gupta
	fmt.Println("=== 1. Purchasing Ticket ===")
	r, err := client.PurchaseTicket(ctx, &pb.PurchaseRequest{
		User: &pb.User{
			FirstName: "Ravi",
			LastName:  "Gupta",
			Email:     "Ravi.Gupta@example.com",
		},
	})
	if err != nil {
		log.Fatalf("Ticket purchase failed: %v", err)
	}
	fmt.Printf("Ticket purchased successfully!\nDetails:\n  From: %s\n  To: %s\n  Name: %s %s\n  Seat: %s\n  Price: $%.2f\n\n",
		r.Ticket.From, r.Ticket.To, r.Ticket.User.FirstName, r.Ticket.User.LastName, r.Ticket.Seat, r.Ticket.Price)

	// 2. Get the receipt for Ravi Gupta
	fmt.Println("=== 2. Retrieving Receipt ===")
	receipt, err := client.GetReceipt(ctx, &pb.ReceiptRequest{Email: "Ravi.Gupta@example.com"})
	if err != nil {
		log.Fatalf("Failed to retrieve receipt: %v", err)
	}
	fmt.Printf("Receipt retrieved:\n  From: %s\n  To: %s\n  Name: %s %s\n  Seat: %s\n  Price: $%.2f\n\n",
		receipt.Ticket.From, receipt.Ticket.To, receipt.Ticket.User.FirstName, receipt.Ticket.User.LastName, receipt.Ticket.Seat, receipt.Ticket.Price)

	// 3. View users in section B
	fmt.Println("=== 3. Viewing Users in Section B ===")
	sectionB, err := client.GetSectionUsers(ctx, &pb.SectionRequest{Section: "B"})
	if err != nil {
		log.Fatalf("Failed to get section B users: %v", err)
	}
	fmt.Printf("Users in Section B:\n")
	for _, ticket := range sectionB.Tickets {
		fmt.Printf("  Name: %s %s, Seat: %s\n", ticket.User.FirstName, ticket.User.LastName, ticket.Seat)
	}
	fmt.Println()

	// 4. View users in section A
	fmt.Println("=== 4. Viewing Users in Section A ===")
	sectionA, err := client.GetSectionUsers(ctx, &pb.SectionRequest{Section: "A"})
	if err != nil {
		log.Fatalf("Failed to get section A users: %v", err)
	}
	fmt.Printf("Users in Section A:\n")
	for _, ticket := range sectionA.Tickets {
		fmt.Printf("  Name: %s %s, Seat: %s\n", ticket.User.FirstName, ticket.User.LastName, ticket.Seat)
	}
	fmt.Println()

	// 5. Modify Ravi's seat to section B
	fmt.Println("=== 5. Modifying Seat ===")
	modifiedSeat, err := client.ModifySeat(ctx, &pb.ModifySeatRequest{
		Email:   "Ravi.Gupta@example.com",
		NewSeat: "B",
	})
	if err != nil {
		log.Fatalf("Failed to modify seat: %v", err)
	}
	fmt.Printf("Seat modified successfully!\nNew Details:\n  Name: %s %s\n  New Seat: %s\n\n",
		modifiedSeat.Ticket.User.FirstName, modifiedSeat.Ticket.User.LastName, modifiedSeat.Ticket.Seat)

	// 6.a Verify users in section B
	fmt.Println("=== 6.a Verifying Users in Section B ===")
	sectionB, err = client.GetSectionUsers(ctx, &pb.SectionRequest{Section: "B"})
	if err != nil {
		log.Fatalf("Failed to get section B users: %v", err)
	}
	fmt.Printf("Users in Section B:\n")
	for _, ticket := range sectionB.Tickets {
		fmt.Printf("  Name: %s %s, Seat: %s\n", ticket.User.FirstName, ticket.User.LastName, ticket.Seat)
	}
	fmt.Println()

	// 6.b Verify users in section A
	fmt.Println("=== 6.b Verifying Users in Section A ===")
	sectionA, err = client.GetSectionUsers(ctx, &pb.SectionRequest{Section: "A"})
	if err != nil {
		log.Fatalf("Failed to get section A users: %v", err)
	}
	fmt.Printf("Users in Section A:\n")
	for _, ticket := range sectionA.Tickets {
		fmt.Printf("  Name: %s %s, Seat: %s\n", ticket.User.FirstName, ticket.User.LastName, ticket.Seat)
	}
	fmt.Println()

	// 7. Remove User Ravi Gupta from the train
	fmt.Println("=== 7. Removing User ===")
	removeRes, err := client.RemoveUser(ctx, &pb.RemoveUserRequest{Email: "Ravi.Gupta@example.com"})
	if err != nil {
		log.Fatalf("Failed to remove user: %v", err)
	}
	if removeRes.Success {
		fmt.Println("User Ravi.Gupta@example.com removed successfully")
	} else {
		fmt.Println("Failed to remove user Ravi.Gupta@example.com.")
	}

	// 8. Attempt to get the receipt after removal to verify deletion
	fmt.Println("\n=== 8. Verifying Deletion ===")
	_, err = client.GetReceipt(ctx, &pb.ReceiptRequest{Email: "Ravi.Gupta@example.com"})
	if err != nil {
		// Extract the specific error message from the full gRPC error
		if grpcErr, ok := status.FromError(err); ok {
			fmt.Println(grpcErr.Message()) // Print only the message part of the gRPC error
		} else {
			fmt.Println(err.Error())
		}
	} else {
		log.Fatalf("Unexpectedly found the receipt for a removed user!")
	}
}
