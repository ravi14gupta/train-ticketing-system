package ticket

import (
	"context"
	"fmt"
	"sync"

	pb "github.com/ravi14gupta/train-ticketing-system/proto"
)

type Service struct {
	pb.UnimplementedTicketServiceServer
	ticketsByEmail map[string]*pb.Ticket
	sectionSeats   map[string][]*pb.Ticket
	mutex          sync.Mutex
}

func NewService() *Service {
	return &Service{
		ticketsByEmail: make(map[string]*pb.Ticket),
		sectionSeats:   make(map[string][]*pb.Ticket),
	}
}

func (s *Service) PurchaseTicket(ctx context.Context, req *pb.PurchaseRequest) (*pb.PurchaseResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Determine the seat section based on the round-robin strategy
	seat := "A"
	if len(s.sectionSeats["A"]) > len(s.sectionSeats["B"]) {
		seat = "B"
	}

	ticket := &pb.Ticket{
		From:  "London",
		To:    "France",
		User:  req.User,
		Price: 20.00,
		Seat:  seat,
	}

	// Replace any existing ticket for the user
	if oldTicket, exists := s.ticketsByEmail[req.User.Email]; exists {
		// Remove the old ticket from the sectionSeats map
		oldSeat := oldTicket.Seat
		sectionTickets := s.sectionSeats[oldSeat]
		for i, t := range sectionTickets {
			if t.User.Email == req.User.Email {
				s.sectionSeats[oldSeat] = append(sectionTickets[:i], sectionTickets[i+1:]...)
				break
			}
		}
	}

	// Store the new ticket
	s.ticketsByEmail[req.User.Email] = ticket
	s.sectionSeats[seat] = append(s.sectionSeats[seat], ticket)

	return &pb.PurchaseResponse{Ticket: ticket}, nil
}

func (s *Service) GetReceipt(ctx context.Context, req *pb.ReceiptRequest) (*pb.ReceiptResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Retrieve the ticket for the given email
	ticket, exists := s.ticketsByEmail[req.Email]
	if !exists {
		return nil, fmt.Errorf("ticket not found for email: %s", req.Email)
	}
	return &pb.ReceiptResponse{Ticket: ticket}, nil
}

func (s *Service) GetSectionUsers(ctx context.Context, req *pb.SectionRequest) (*pb.SectionResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Retrieve users by section
	sectionTickets, exists := s.sectionSeats[req.Section]
	if !exists {
		return &pb.SectionResponse{Tickets: []*pb.Ticket{}}, nil
	}
	return &pb.SectionResponse{Tickets: sectionTickets}, nil
}

func (s *Service) RemoveUser(ctx context.Context, req *pb.RemoveUserRequest) (*pb.RemoveUserResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Remove the user's ticket
	ticket, exists := s.ticketsByEmail[req.Email]
	if !exists {
		return &pb.RemoveUserResponse{Success: false}, fmt.Errorf("user not found with email: %s", req.Email)
	}

	// Remove the ticket from the sectionSeats map
	seat := ticket.Seat
	sectionTickets := s.sectionSeats[seat]
	for i, t := range sectionTickets {
		if t.User.Email == req.Email {
			s.sectionSeats[seat] = append(sectionTickets[:i], sectionTickets[i+1:]...)
			break
		}
	}

	delete(s.ticketsByEmail, req.Email)
	return &pb.RemoveUserResponse{Success: true}, nil
}

func (s *Service) ModifySeat(ctx context.Context, req *pb.ModifySeatRequest) (*pb.ModifySeatResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Find the user's ticket
	ticket, exists := s.ticketsByEmail[req.Email]
	if !exists {
		return nil, fmt.Errorf("user not found with email: %s", req.Email)
	}

	// Remove the ticket from the old section
	oldSeat := ticket.Seat
	sectionTickets := s.sectionSeats[oldSeat]
	for i, t := range sectionTickets {
		if t.User.Email == req.Email {
			s.sectionSeats[oldSeat] = append(sectionTickets[:i], sectionTickets[i+1:]...)
			break
		}
	}

	// Update the seat and add the ticket to the new section
	ticket.Seat = req.NewSeat
	s.sectionSeats[req.NewSeat] = append(s.sectionSeats[req.NewSeat], ticket)

	return &pb.ModifySeatResponse{Ticket: ticket}, nil
}
