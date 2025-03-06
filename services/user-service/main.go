package main

import (
	"log"
	"net"

	"github.com/Aditya-PS-05/NeetChamp/user-service/controllers"
	"github.com/Aditya-PS-05/NeetChamp/user-service/database"

	pb "github.com/Aditya-PS-05/NeetChamp/shared-libs/proto"
	"google.golang.org/grpc"
)

func main() {
	// Connect to the database
	database.ConnectDatabase()

	// Start gRPC server
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userController := &controllers.UserController{DB: database.DB}

	pb.RegisterUserServiceServer(grpcServer, userController)

	log.Println("✅ User service is running on port 50052...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
