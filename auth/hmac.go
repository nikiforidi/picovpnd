package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var SharedSecret = []byte("supersecretkey")

func GenerateHMAC(message string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

// --- Server Interceptor for HMAC Auth ---
func HMACAuthInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}
	ts := md.Get("x-auth-timestamp")
	sig := md.Get("x-auth-signature")
	if len(ts) == 0 || len(sig) == 0 {
		return nil, fmt.Errorf("missing auth headers")
	}
	expectedSig := GenerateHMAC(ts[0], SharedSecret)
	if !hmac.Equal([]byte(sig[0]), []byte(expectedSig)) {
		return nil, fmt.Errorf("invalid signature")
	}
	return handler(ctx, req)
}

// --- Client with HMAC Auth ---
// func grpcClient() {
// 	conn, err := grpc.NewClient(
// 		"localhost:50051",
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 		// grpc.WithBlock(),
// 	)
// 	if err != nil {
// 		log.Fatalf("did not connect: %v", err)
// 	}
// 	defer conn.Close()
// 	c := pb.NewGreeterClient(conn)

// 	ts := fmt.Sprintf("%d", time.Now().Unix())
// 	sig := GenerateHMAC(ts, sharedSecret)
// 	md := metadata.Pairs(
// 		"x-auth-timestamp", ts,
// 		"x-auth-signature", sig,
// 	)
// 	ctx := metadata.NewOutgoingContext(context.Background(), md)
// 	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "world"})
// 	if err != nil {
// 		log.Fatalf("could not greet: %v", err)
// 	}
// 	log.Printf("Greeting: %s", r.GetMessage())
// }

// func main() {
// 	// Start gRPC server
// 	go func() {
// 		lis, err := net.Listen("tcp", ":50051")
// 		if err != nil {
// 			log.Fatalf("failed to listen: %v", err)
// 		}
// 		s := grpc.NewServer(grpc.UnaryInterceptor(hmacAuthInterceptor))
// 		pb.RegisterGreeterServer(s, &server{})
// 		log.Println("gRPC server listening on :50051")
// 		if err := s.Serve(lis); err != nil {
// 			log.Fatalf("failed to serve: %v", err)
// 		}
// 	}()
// 	time.Sleep(1 * time.Second) // Wait for server to start
// 	grpcClient()
// }
