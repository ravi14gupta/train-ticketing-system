# Train Ticket Booking System

## Overview
This project is a train ticket booking system implemented using **Golang** and **gRPC**. The system simulates the process of purchasing a train ticket, retrieving receipt details, viewing seat allocations, removing users from the train, and modifying user seat assignments. The data is stored in memory, and all results are outputted to the console.

## Features
- **Purchase Ticket**: Submit a purchase for a train ticket from London to France. The ticket costs $20 and includes the user's details (first name, last name, email).
- **View Receipt**: Retrieve the details of the receipt for a specific user, including the journey details, user information, price paid, and seat allocated.
- **View Section Users**: View all users and their allocated seats within a specific section (Section A or Section B).
- **Remove User**: Remove a user from the train, freeing up their allocated seat.
- **Modify Seat**: Modify the seat allocation for a specific user.

## Prerequisites
Before you begin, ensure you have the following:

**1. Host Operating System**
Ubuntu 20.04.6 LTS or Any Linux distribution OS

**2. Go Programming Language**
Version 1.23.0 or later 
- Download and Install From go dev ( https://go.dev/doc/install )
  ```bash
  sudo rm -rf /usr/local/go
  wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
  sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz
  export PATH=$PATH:/usr/local/go/bin
  go version
  ```

**3. gRPC Tools:**
- protoc compiler
  ```bash
  sudo apt-get install -y protobuf-compiler
  protocol --version
  ```
- protoc-gen-go plugin
  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  ```
- protoc-gen-go-grpc plugin
  ```bash
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  ```

**4. Network Access:**
Ensure network access to install required packages and dependencies.

## Getting Started
To get a copy of the project up and running on your local machine, follow these steps:

**1. Clone the Repository**

Clone the repository to your local machine using the following command:
```bash
git clone https://github.com/ravi14gupta/train-ticketing-system.git
```

**2. Navigate to the Project Directory**

Change into the project directory:
```bash
cd train-ticketing-system
```

**3. Install Dependencies**

Ensure all dependencies are installed by running:
```bash
go mod tidy
```
This will download the necessary Go modules, including gRPC.

**4. [OPTIONAL] Compile Protocol Buffers**

If you make any changes to the .proto files, compile them using:
```bash
protoc --go_out=. --go-grpc_out=. proto/ticket.proto
```

**5. Run the gRPC Server**

Start the gRPC server by running:
```bash
go run cmd/server/main.go
```

**6. Run the gRPC Client**

In a separate terminal, run the gRPC client to interact with the server:
```bash
go run cmd/client/main.go
```

**7. Testing**

To run the unit tests, use the following command:
```bash
go test ./...
```
## Limitations
- **Multiple Tickets For Same User:**
  Currently, it overwrites the existing ticket, if a user with the same email purchases a ticket multiple times.

## Future Enhancements
Here are some potential future enhancements that could be implemented:

- **Persistent Storage**
Integrate a database to store ticket and user data persistently, instead of using in-memory storage.

- **Advanced Seat Allocation**
Develop a more sophisticated seat allocation algorithm that considers user preferences, such as window or aisle seats, and optimizes seat distribution across sections.

- **Payment Integration**
Integrate with a payment gateway.

- **Security Enhancements**
Add authentication and authorization mechanisms, to secure the APIs. Implement role-based access control (RBAC) for different user roles (e.g., admin, customer).

- **Multi-Train Support**
Extend the system to handle multiple trains, allowing users to choose their preferred train for a specific route and time.

- **Load Testing and Performance Optimization**
Implement load testing to simulate high traffic and optimize the system's performance to handle large numbers of concurrent users.

- **Web or Mobile Interface**
Develop a user-friendly web or mobile interface to interact with the gRPC APIs, making the system more accessible to end-users.
