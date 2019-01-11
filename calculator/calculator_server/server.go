package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/4nc3str4l/GoLangApi_MicroservicesCourse/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

// Sum implementation
func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", req)
	p1 := req.GetParameters().GetP1()
	p2 := req.GetParameters().GetP2()
	result := p1 + p2
	res := &calculatorpb.SumResponse{
		Result: result,
	}
	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("PrimeNumberDecomposition function was invoked with %v\n", req)
	num := req.GetNum()
	k := int32(2)
	for num > 1 {
		// Found a factor
		if num%k == 0 {
			res := &calculatorpb.PrimeResponse{
				Result: k,
			}
			stream.Send(res)
			num = num / k
			time.Sleep(1000 * time.Millisecond)
		} else {
			k++
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("ComputeAverage function was invoked with %v\n", stream)
	sum := 0.0
	reqLength := 0.0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// We have finished reading the client stream
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Result: sum / reqLength,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		sum += float64(req.GetNum())
		reqLength++
	}
}

func main() {

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	// Register the calculator service
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
