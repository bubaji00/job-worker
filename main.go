// go build -o timer.exe
// .\timer.exe start --time=[int] --unit=[sec/min/hr]

//since last meeting: adjusted mutex lock and unlock, refactoring into different file, OOD, prevent large time input (<100hr)

// at the end of the project, user can start a timer and when the time is reached it will execute some cmd command on a remote machine

package main

import (
	"fmt"
	"google.golang.org/grpc"
	pb "jobWorker/proto"
	"jobWorker/server"
	"jobWorker/worker"
	"log"
	"net"
)

func main() {
	var jobDispatch = worker.JobDispatcher{
		Jobs: make(map[string]*worker.Job),
	}

	// 2. Listen on a port for gRPC
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on :50051: %v", err)
	}

	// 3. Create a gRPC server
	grpcServer := grpc.NewServer()

	// 4. Create and register your JobManagerServer
	jobManagerServer := server.NewJobManagerServer(&jobDispatch)
	pb.RegisterJobManagerServer(grpcServer, jobManagerServer)

	fmt.Println("gRPC server running on port 50051")

	// 5. Start serving
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}

	/*
		rootCmd := cmd.SetupCommands(&jobDispatch)
		if err := rootCmd.Execute(); err != nil {
			fmt.Println(err)
		}
		jobDispatch.Start()
	*/
}
