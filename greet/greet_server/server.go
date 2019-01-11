package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/4nc3str4l/GoLangApi_MicroservicesCourse/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

// Greet implementation
func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

// Sum implementation
func (*server) Sum(ctx context.Context, req *greetpb.SumRequest) (*greetpb.SumResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", req)
	p1 := req.GetParameters().GetP1()
	p2 := req.GetParameters().GetP2()
	result := p1 + p2
	res := &greetpb.SumResponse{
		Result: result,
	}
	return res, nil
}

func main() {

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	// Register the calculator service
	greetpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
