package main

import (
	"Agri/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterSupplierServiceProtoServer(grpcServer, &proto.SupplierServiceProto{})

	log.Println("🚀 gRPC Server running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
