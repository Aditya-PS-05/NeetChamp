package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Aditya-PS-05/NeetChamp/auth-service/controllers"
	"github.com/Aditya-PS-05/NeetChamp/auth-service/database"
	"github.com/Aditya-PS-05/NeetChamp/auth-service/utils"
	"github.com/Aditya-PS-05/NeetChamp/shared-libs/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// 🏆 Connect to the database (optimized with connection pooling)
	database.ConnectDatabase()

	// 🏆 Connect to Redis (efficient token storage)
	utils.ConnectRedis()

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("❌ Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.MaxConcurrentStreams(2000), // Allow 200 parallel requests
	)
	proto.RegisterAuthServiceServer(grpcServer, &controllers.AuthServiceServer{})

	// 🔹 Enable gRPC reflection
	reflection.Register(grpcServer)

	fmt.Println("🚀 Auth Service is running on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("❌ Failed to serve: %v", err)
	}
}
