package main

func main() {
	// creds, err := credentials.NewClientTLSFromFile("../cert.pem", "service1")
	// if err != nil {
	// 	log.Fatalf("failed to load TLS cert: %v", err)
	// }
	// conn, err := grpc.NewClient("service1:50051", grpc.WithTransportCredentials(creds))
	// if err != nil {
	// 	log.Fatalf("did not connect: %v", err)
	// }
	// defer conn.Close()
	// c := pb.NewOpenConnectServiceClient(conn)
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()
	// r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "World"})
	// if err != nil {
	// 	log.Fatalf("could not greet: %v", err)
	// }
	// log.Printf("Greeting: %s", r.Message)
}
