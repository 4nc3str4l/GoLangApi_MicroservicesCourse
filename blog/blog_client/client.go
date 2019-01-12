package main

import (
	"context"
	"fmt"
	"log"

	"github.com/4nc3str4l/GoLangApi_MicroservicesCourse/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog Client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	fmt.Println("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "Cristian",
		Title:    "My First Blog",
		Content:  "Content of the first blog",
	}
	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Println("Blog has been created: %v", createBlogRes)
	blogID := createBlogRes.GetBlog().GetId()

	// Read blog
	_, err2 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "sadasd"})
	if err2 != nil {
		fmt.Printf("Error happened while reading: %v\n", err2)
	}

	readblogReq := &blogpb.ReadBlogRequest{BlogId: blogID}
	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readblogReq)
	if readBlogErr != nil {
		fmt.Printf("Error happened while reading; %v", readBlogErr)
	}
	fmt.Printf("Blog was read: %v", readBlogRes)
}
