package grpc

import (
	"context"
	"log"
	"net"

	"github.com/anatolio-deb/picovpnd/ocserv"
	pb "github.com/anatolio-deb/picovpnd/user" // Adjust import path to your generated proto package

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserServiceServer
}

func (s *server) AddUser(ctx context.Context, req *pb.AddUserRequest) (*pb.UserResponse, error) {
	err := ocserv.UserAdd(req.Username, req.Password)
	if err != nil {
		return &pb.UserResponse{Code: 1, Error: err.Error()}, nil
	}
	return &pb.UserResponse{Code: 0, Error: ""}, nil
}

func (s *server) LockUser(ctx context.Context, req *pb.LockUserRequest) (*pb.UserResponse, error) {
	err := ocserv.UserLock(req.Username)
	if err != nil {
		return &pb.UserResponse{Code: 1, Error: err.Error()}, nil
	}
	return &pb.UserResponse{Code: 0, Error: ""}, nil
}

func (s *server) UnlockUser(ctx context.Context, req *pb.LockUserRequest) (*pb.UserResponse, error) {
	err := ocserv.UserUnlock(req.Username)
	if err != nil {
		return &pb.UserResponse{Code: 1, Error: err.Error()}, nil
	}
	return &pb.UserResponse{Code: 0, Error: ""}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &server{})
	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
