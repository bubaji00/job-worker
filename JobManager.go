package server

import (
	"context"
	pb "jobWorker/proto"
	"jobWorker/worker"
	"time"
)

type JobManagerServer struct {
	pb.UnimplementedJobManagerServer
	dispatcher *worker.JobDispatcher
}

func NewJobManagerServer(dispatcher *worker.JobDispatcher) *JobManagerServer {
	return &JobManagerServer{dispatcher: dispatcher}
}
func (s *JobManagerServer) Start(ctx context.Context, req *pb.JobCreationRequest) (*pb.Job, error) {
	duration := 5 * time.Second // test time

	jobID, err := s.dispatcher.StartTimerCore(duration, req.Name)
	if err != nil {
		return nil, err
	}

	return &pb.Job{
		Id:   jobID,
		Name: req.Name,
	}, nil
}

func (s *JobManagerServer) Stop(ctx context.Context, req *pb.JobID) (*pb.Status, error) {
	err := s.dispatcher.StopTimerCore(req.Id)
	if err != nil {
		return &pb.Status{Status: "Error", ErrorMessage: err.Error()}, nil
	}
	return &pb.Status{Status: "Stopped", JobId: req.Id}, nil
}

func (s *JobManagerServer) Query(ctx context.Context, req *pb.JobID) (*pb.Status, error) {
	state, err := s.dispatcher.QueryTimerCore(req.Id)
	if err != nil {
		return &pb.Status{Status: "Error", ErrorMessage: err.Error()}, nil
	}
	return &pb.Status{Status: state, JobId: req.Id}, nil
}

func (s *JobManagerServer) List(ctx context.Context, req *pb.NilMessage) (*pb.JobStatusList, error) {
	s.dispatcher.Mutex.Lock()
	defer s.dispatcher.Mutex.Unlock()

	var statuses []*pb.Status
	for id, job := range s.dispatcher.Jobs {
		statuses = append(statuses, &pb.Status{
			Status:    job.State,
			JobId:     id,
			IsRunning: job.State == "STARTED",
		})
	}
	return &pb.JobStatusList{JobStatusList: statuses}, nil
}
