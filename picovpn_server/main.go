package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/Netflix/go-expect"
	pb "github.com/anatolio-deb/picovpnd/picovpnd"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.OpenConnectServiceServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) UserAdd(_ context.Context, req *pb.UserAddRequest) (*pb.Response, error) {
	c, err := expect.NewConsole(expect.WithStdout(os.Stdout))
	if err != nil {
		return nil, err
	}
	defer c.Close()

	cmd := exec.Command("ocpasswd", "-c", "/etc/ocserv/passwd", req.Username)
	cmd.Stdin = c.Tty()
	cmd.Stdout = c.Tty()
	cmd.Stderr = c.Tty()

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	go func() {
		c.ExpectString("Enter password:")
	}()

	time.Sleep(time.Second)
	c.SendLine(req.Password)

	go func() {
		c.ExpectString("Re-enter password:")
	}()

	time.Sleep(time.Second)
	c.SendLine(req.Password)

	err = cmd.Wait()
	return &pb.Response{
		Error: err.Error(),
	}, err
}

func (s *server) UserLock(context.Context, *pb.UserLockRequest) (*pb.Response, error) {
	return &pb.Response{
		Error: "Not implemented",
	}, fmt.Errorf("not implemented")

}

func (s *server) UserUnlock(context.Context, *pb.UserUnlockRequest) (*pb.Response, error) {
	return &pb.Response{
		Error: "Not implemented",
	}, fmt.Errorf("not implemented")
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

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterOpenConnectServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
