package ticket

import (
	"context"
	"testing"

	pb "github.com/ravi14gupta/train-ticketing-system/proto"
	"github.com/stretchr/testify/assert"
)

func TestPurchaseTicket(t *testing.T) {
	service := NewService()

	user1 := &pb.User{
		FirstName: "Raju",
		LastName:  "Gupta",
		Email:     "Raju.Gupta@example.com",
	}

	user2 := &pb.User{
		FirstName: "Ravi",
		LastName:  "Gupta",
		Email:     "Ravi.Gupta@example.com",
	}

	// Purchase tickets for two users
	res1, err := service.PurchaseTicket(context.Background(), &pb.PurchaseRequest{User: user1})
	assert.NoError(t, err)
	assert.Equal(t, "Raju", res1.Ticket.User.FirstName)
	assert.Equal(t, "A", res1.Ticket.Seat)

	res2, err := service.PurchaseTicket(context.Background(), &pb.PurchaseRequest{User: user2})
	assert.NoError(t, err)
	assert.Equal(t, "Ravi", res2.Ticket.User.FirstName)
	assert.Equal(t, "B", res2.Ticket.Seat)

	// Check the internal state
	assert.Equal(t, 1, len(service.sectionSeats["A"]))
	assert.Equal(t, 1, len(service.sectionSeats["B"]))
	assert.Equal(t, user1.Email, service.sectionSeats["A"][0].User.Email)
	assert.Equal(t, user2.Email, service.sectionSeats["B"][0].User.Email)
}

func TestGetReceipt(t *testing.T) {
	service := NewService()

	user := &pb.User{
		FirstName: "Raju",
		LastName:  "Gupta",
		Email:     "Raju.Gupta@example.com",
	}

	_, _ = service.PurchaseTicket(context.Background(), &pb.PurchaseRequest{User: user})
	res, err := service.GetReceipt(context.Background(), &pb.ReceiptRequest{Email: user.Email})

	assert.NoError(t, err)
	assert.Equal(t, "Raju", res.Ticket.User.FirstName)
	assert.Equal(t, "London", res.Ticket.From)
	assert.Equal(t, "France", res.Ticket.To)
	assert.Equal(t, float64(20.00), res.Ticket.Price)
	assert.Equal(t, "A", res.Ticket.Seat)
}

func TestGetSectionUsers(t *testing.T) {
	service := NewService()

	user1 := &pb.User{
		FirstName: "Raju",
		LastName:  "Gupta",
		Email:     "Raju.Gupta@example.com",
	}

	user2 := &pb.User{
		FirstName: "Ravi",
		LastName:  "Gupta",
		Email:     "Ravi.Gupta@example.com",
	}

	_, _ = service.PurchaseTicket(context.Background(), &pb.PurchaseRequest{User: user1})
	_, _ = service.PurchaseTicket(context.Background(), &pb.PurchaseRequest{User: user2})

	// Get users in Section A
	resA, err := service.GetSectionUsers(context.Background(), &pb.SectionRequest{Section: "A"})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resA.Tickets))
	assert.Equal(t, "Raju", resA.Tickets[0].User.FirstName)

	// Get users in Section B
	resB, err := service.GetSectionUsers(context.Background(), &pb.SectionRequest{Section: "B"})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(resB.Tickets))
	assert.Equal(t, "Ravi", resB.Tickets[0].User.FirstName)
}

func TestRemoveUser(t *testing.T) {
	service := NewService()

	user := &pb.User{
		FirstName: "Raju",
		LastName:  "Gupta",
		Email:     "Raju.Gupta@example.com",
	}

	_, _ = service.PurchaseTicket(context.Background(), &pb.PurchaseRequest{User: user})

	// Remove the user
	res, err := service.RemoveUser(context.Background(), &pb.RemoveUserRequest{Email: user.Email})
	assert.NoError(t, err)
	assert.True(t, res.Success)

	// Try to get receipt after removal, should fail
	_, err = service.GetReceipt(context.Background(), &pb.ReceiptRequest{Email: user.Email})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ticket not found for email")

	// Check the internal state
	assert.Equal(t, 0, len(service.sectionSeats["A"]))
	assert.Nil(t, service.ticketsByEmail[user.Email])
}

func TestModifySeat(t *testing.T) {
	service := NewService()

	user := &pb.User{
		FirstName: "Raju",
		LastName:  "Gupta",
		Email:     "Raju.Gupta@example.com",
	}

	_, _ = service.PurchaseTicket(context.Background(), &pb.PurchaseRequest{User: user})

	// Modify the user's seat
	modRes, err := service.ModifySeat(context.Background(), &pb.ModifySeatRequest{
		Email:   user.Email,
		NewSeat: "B",
	})
	assert.NoError(t, err)
	assert.Equal(t, "B", modRes.Ticket.Seat)

	// Check the internal state
	assert.Equal(t, 0, len(service.sectionSeats["A"]))
	assert.Equal(t, 1, len(service.sectionSeats["B"]))
	assert.Equal(t, "Raju", service.sectionSeats["B"][0].User.FirstName)
}

func TestPurchaseTicketSeatAllocation(t *testing.T) {
	service := NewService()

	// Create multiple users to fill both sections
	users := []*pb.User{
		{FirstName: "User1", LastName: "One", Email: "user1@example.com"},
		{FirstName: "User2", LastName: "Two", Email: "user2@example.com"},
		{FirstName: "User3", LastName: "Three", Email: "user3@example.com"},
		{FirstName: "User4", LastName: "Four", Email: "user4@example.com"},
	}

	for i, user := range users {
		res, err := service.PurchaseTicket(context.Background(), &pb.PurchaseRequest{User: user})
		assert.NoError(t, err)
		expectedSeat := "A"
		if i%2 != 0 {
			expectedSeat = "B"
		}
		assert.Equal(t, expectedSeat, res.Ticket.Seat)
	}

	// Verify seat distribution
	assert.Equal(t, 2, len(service.sectionSeats["A"]))
	assert.Equal(t, 2, len(service.sectionSeats["B"]))
}
