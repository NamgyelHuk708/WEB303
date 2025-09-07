package main

import (
	"context"
	"log"
	"net"
	pb "practical-one/proto/gen/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServiceServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Received request for greeting: %s", in.Name)
	message := "Hello, " + in.Name + "!"
	return &pb.HelloResponse{Message: message}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServiceServer(s, &server{})
	log.Printf("Greeter service listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}