package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/Netflix/go-expect"
	pb "github.com/anatolio-deb/picovpnd/picovpnd"
	"golang.org/x/crypto/acme/autocert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

func (s *server) UserLock(_ context.Context, req *pb.UserLockRequest) (*pb.Response, error) {
	b, err := exec.Command("ocpasswd", "--lock", req.Username).CombinedOutput()
	return &pb.Response{
		Error: string(b),
	}, err
}

func (s *server) UserUnlock(_ context.Context, req *pb.UserUnlockRequest) (*pb.Response, error) {
	b, err := exec.Command("ocpasswd", "--unlock", req.Username).CombinedOutput()
	return &pb.Response{
		Error: string(b),
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

func main() {
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(os.Getenv("AUTOCERT_DOMAIN")), // Use your email or domain here
		Cache:      autocert.DirCache(os.Getenv("AUTOCERT_DIR")),         // Directory to cache certificates
	}

	flag.Parse()

	lis, err := net.Listen("tcp", ":0") // Listen on a random port
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create tls based credential.
	cert, err := certManager.GetCertificate(&tls.ClientHelloInfo{
		ServerName: os.Getenv("AUTOCERT_DOMAIN"),
	})

	if err != nil {
		log.Fatalf("failed to get certificate: %v", err)
	}
	// Create credentials from the certificate.
	creds := credentials.NewServerTLSFromCert(cert)
	if err != nil {
		log.Fatalf("failed to create credentials: %v", err)
	}

	s := grpc.NewServer(grpc.Creds(creds))

	// Register EchoServer on the server.
	pb.RegisterOpenConnectServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
