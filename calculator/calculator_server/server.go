package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"time"

	"github.com/4nc3str4l/GoLangApi_MicroservicesCourse/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Printf("FindMaximum function was invoked with %v\n", stream)

	// Get the minimum possible integer
	max := int64(-(int(^uint(0)>>1) - 1))

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Eror while reading client stream: %v", err)
			return err
		}

		currentNum := req.GetNum()
		if max < currentNum {
			max = currentNum
			sendErr := stream.Send(&calculatorpb.FindMaximumResponse{
				Result: max,
			})
			if sendErr != nil {
				log.Fatalf("Eror while sending data to client stream: %v", sendErr)
				return sendErr
			}
		}

	}
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Println("Reveived SquareRoot RPC")
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %v", number),
		)
	}

	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
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
