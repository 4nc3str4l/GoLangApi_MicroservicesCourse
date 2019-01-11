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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	//doClientStreaming(c)
	//doBiDiStreaming(c)
	doErorUnary(c)
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

func doBiDiStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Bidi Streaming RPC...")

	// Create a stream by invoking the client
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	// Create a wait channel to block on it
	waitc := make(chan struct{})

	go func() {
		numbers := []int64{1, 5, 3, 6, 2, 20}

		// Send a bunch of messages to the client (go routine)
		for _, num := range numbers {
			fmt.Printf("Sending message: %v\n", num)
			stream.Send(&calculatorpb.FindMaximumRequest{
				Num: num,
			})
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

func doErorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a SquareRoot Unary RPC...")
	// Correct call
	doErrorCall(c, 10)

	// Wrong call
	doErrorCall(c, -1)
}

func doErrorCall(c calculatorpb.CalculatorServiceClient, n int32) {

	// Correct
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: n})

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			// Actual error from gRPC (user Error)
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number!")
				return
			}
		} else {
			log.Fatalf("Big error calling SquareRoot: %v", err)
			return
		}
	}
	fmt.Printf("Result of square root of %v is %v\n", n, res.GetNumberRoot())
}
