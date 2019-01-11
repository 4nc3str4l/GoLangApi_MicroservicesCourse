package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/4nc3str4l/GoLangApi_MicroservicesCourse/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	doBiDiStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Cristian",
			LastName:  "Muriel",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streamign RPC...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Cristian",
			LastName:  "Muriel",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GreetingManyTimes RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// End of stream reached
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream %v", err)
		}
		log.Printf("Response from GreetManyTimes %v", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Cristian",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Myriam",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Guillem",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Pierre",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "David",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling longGreet: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending request %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("error while receiving response from LongGreet: %v", err)
	}
	fmt.Printf("LongGreet Response: %v\n", res)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Bidi Streaming RPC...")

	// Create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Cristian",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Myriam",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Guillem",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Pierre",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "David",
			},
		},
	}

	// Create a wait channel to block on it
	waitc := make(chan struct{})

	go func() {
		// Send a bunch of messages to the client (go routine)
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		// Receive a bunch of messages from the client (go routine)
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				// Close the wait channel
				close(waitc)
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				close(waitc)
			}
			fmt.Printf("Receved: %v\n", res.GetResult())
		}

	}()

	// block until everything is done
	<-waitc
}
