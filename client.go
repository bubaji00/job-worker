package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb "jobWorker/proto"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewJobManagerClient(conn)

	// start job
	job, err := client.Start(context.Background(), &pb.JobCreationRequest{
		Name: "FirstTimer",
	})
	if err != nil {
		log.Fatalf("Could not start job: %v", err)
	}
	fmt.Printf("Started job with ID: %s\n", job.Id)

	// query
	statusResp, err := client.Query(context.Background(), &pb.JobID{Id: job.Id})
	if err != nil {
		log.Fatalf("Could not query job: %v", err)
	}
	fmt.Printf("Job status: %s\n", statusResp.Status)

	//stop
	stopResp, err := client.Stop(context.Background(), &pb.JobID{Id: job.Id})
	if err != nil {
		log.Fatalf("Could not stop job: %v", err)
	}
	fmt.Printf("Stopped job. New status: %s\n", stopResp.Status)
}
