package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/anatolio-deb/picovpnd/api"
	"github.com/anatolio-deb/picovpnd/auth"
	"github.com/anatolio-deb/picovpnd/core"
	pb "github.com/anatolio-deb/picovpnd/grpc"
	"github.com/anatolio-deb/picovpnd/ip"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.OpenConnectServiceServer
}

func (s *server) UserAdd(_ context.Context, req *pb.UserAddRequest) (*pb.Response, error) {
	err := core.UserAdd(req.Username, req.Password)
	return &pb.Response{
		Error: err.Error(),
	}, err
}

func (s *server) UserLock(_ context.Context, req *pb.UserLockRequest) (*pb.Response, error) {
	err := core.UserLock(req.Username)
	return &pb.Response{
		Error: err.Error(),
	}, err
}

func (s *server) UserUnlock(_ context.Context, req *pb.UserUnlockRequest) (*pb.Response, error) {
	err := core.UserUnlock(req.Username)
	return &pb.Response{
		Error: err.Error(),
	}, err
}

func (s *server) UserDelete(context.Context, *pb.UserDeleteRequest) (*pb.Response, error) {
	return &pb.Response{
		Error: "Not implemented",
	}, fmt.Errorf("not implemented")
}

func (s *server) UserChangePassword(context.Context, *pb.UserChangePasswordRequest) (*pb.Response, error) {
	return &pb.Response{
		Error: "Not implemented",
	}, fmt.Errorf("not implemented")
}

func GetCert(ctx context.Context, req *pb.AuthenticateRequest, opts ...grpc.CallOption) (*pb.CertResponse, error) {
	cert, _, err := auth.NewSSLCertAndKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate cert and key: %v", err)
	}
	certb, err := os.ReadFile(cert.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to read cert: %v", err)
	}
	// keyb, err := os.ReadFile(key.Name())
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to read key: %v", err)
	// }
	return &pb.CertResponse{
		Cert: string(certb),
	}, nil
}

// https://github.com/grpc/grpc-go/blob/master/examples/features/encryption/TLS/server/main.go
func main() {
	daemonPort := os.Getenv("DAEMON_PORT")
	cert, key, err := auth.NewSSLCertAndKey()
	if err != nil {
		log.Fatal(err)
	}
	lis, err := net.Listen("tcp", ":"+daemonPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create tls based credential.
	creds, err := credentials.NewServerTLSFromFile(cert.Name(), key.Name())
	if err != nil {
		log.Fatalf("failed to create credentials: %v", err)
	}

	// s := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(auth.HMACAuthInterceptor))
	s := grpc.NewServer(grpc.Creds(creds))

	// Register EchoServer on the server.
	pb.RegisterOpenConnectServiceServer(s, &server{})

	ip, err := ip.GetPublicIP()
	if err != nil {
		log.Fatalf("failed to get public IP: %v", err)
	}
	b, err := os.ReadFile(cert.Name())
	if err != nil {
		log.Fatalf("failed to read certificate: %v", err)
	}

	daemon := api.Daemon{
		Address:     ip,
		Port:        daemonPort,
		Certificate: string(b),
	}

	go api.RegisterSelf(daemon)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
