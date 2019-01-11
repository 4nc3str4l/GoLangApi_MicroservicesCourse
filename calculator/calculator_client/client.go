package main

import (
	"context"
	"fmt"
	"log"

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

	c2 := calculatorpb.NewCalculatorServiceClient(cc)
	doSum(c2, 10, 3)

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
