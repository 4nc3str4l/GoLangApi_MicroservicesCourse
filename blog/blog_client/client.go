package main

import (
	"context"
	"fmt"
	"io"
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

	// Update
	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Cristian Muriel",
		Title:    "My Second Blog",
		Content:  "Content of the first blog with new cool stuff",
	}
	updateRes, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if updateErr != nil {
		fmt.Printf("Error happeed while updating: %v \n", updateErr)
	}
	fmt.Printf("Blog was updated: %v\n", updateRes)

	// Delete blog
	deleteBlogRes, deleteErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogID})
	if deleteErr != nil {
		fmt.Printf("Error happeed while deleteing: %v \n", deleteErr)
	}
	fmt.Printf("Blog was deleted: %v", deleteBlogRes)

	// List blogs
	resStream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("Error while calling ListBlog RPC: %v", err)
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
		log.Printf("Response from ListBlog %v", msg.GetBlog())
	}
}
