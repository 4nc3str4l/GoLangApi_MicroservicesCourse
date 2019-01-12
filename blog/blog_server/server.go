package main

import (
	"fmt"
	"log"
	"net"

	"github.com/4nc3str4l/GoLangApi_MicroservicesCourse/blog/blogpb"
	"google.golang.org/grpc"
)

type server struct{}

func main() {
	fmt.Println("Blog Service Started")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
