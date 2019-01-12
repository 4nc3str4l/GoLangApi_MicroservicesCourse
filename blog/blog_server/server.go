package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/4nc3str4l/GoLangApi_MicroservicesCourse/blog/blogpb"
	"google.golang.org/grpc"
)

type server struct{}

func main() {
	// if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Blog Service Started")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	// Start the server into a Go Routine
	go func() {
		fmt.Println("Starting Server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to server: %v", err)
		}
	}()

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Blog until the signal is received
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Cloing the listener")
	lis.Close()
	fmt.Println("Shutdown Complete")
}
