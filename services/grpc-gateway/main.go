package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/Aditya-PS-05/NeetChamp/shared-libs/proto/auth"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

const (
	grpcAuthServiceEndpoint = "localhost:50051" // Ensure auth-service is running
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	// Register AuthService to handle REST API requests
	if err := auth.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, grpcAuthServiceEndpoint, opts); err != nil {
		log.Fatalf("❌ Failed to register auth service: %v", err)
	}

	log.Println("🚀 gRPC-Gateway running on port 8080...")
	return http.ListenAndServe(":8080", mux)
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
