package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/4nc3str4l/GoLangApi_MicroservicesCourse/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	//doSum(c, 10, 3)
	//doServerStreaming(c, 120)
	doClientStreaming(c)
}

func doSum(c calculatorpb.CalculatorServiceClient, p1 int32, p2 int32) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &calculatorpb.SumRequest{
		Parameters: &calculatorpb.Parameters{
			P1: p1,
			P2: p2,
		},
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Sum RPC: %v", err)
	}
	log.Printf("Response from Sum: %v", res.Result)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient, num int32) {
	fmt.Println("Starting to do a Server Streaming RPC...")
	toPrint := "Decomposition of " + strconv.Itoa(int(num)) + " is "

	req := &calculatorpb.PrimeRequest{
		Num: num,
	}
	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling PrimeNumberDecomposition RPC: %v", err)
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
		log.Printf("Response from PrimeNumberDecomposition %v", msg.GetResult())
		toPrint = toPrint + strconv.Itoa(int(msg.GetResult())) + "*"
	}

	fmt.Println(toPrint[:len(toPrint)-1])
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	requests := []*calculatorpb.ComputeAverageRequest{
		&calculatorpb.ComputeAverageRequest{
			Num: 1,
		},
		&calculatorpb.ComputeAverageRequest{
			Num: 2,
		},
		&calculatorpb.ComputeAverageRequest{
			Num: 3,
		},
		&calculatorpb.ComputeAverageRequest{
			Num: 4,
		},
	}

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error while calling ComputeAverage: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending request %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("error while receiving response from ComputeAverage: %v", err)
	}
	fmt.Printf("ComputeAverage Response: %v\n", res)
}
