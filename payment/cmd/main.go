package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	paymentAPI "payment/internal/api/payment/v1"
	"payment/internal/config"
	paymentService "payment/internal/service/payment"
	paymentV1 "shared/pkg/proto/payment/v1"
)

const (
	configPath = "./deploy/compose/payment/.env"
)

func main() {
	// Load config
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed load config: %w", err))
	}

	// Start listen
	lis, err := net.Listen("tcp", config.AppConfig().PaymentGRPC.Adress())
	if err != nil {
		log.Printf("Failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("Failed to close listener: %v\n", cerr)
		}
	}()

	// Create GRPC server
	s := grpc.NewServer()

	// Register Service
	service := paymentService.NewService()
	api := paymentAPI.NewAPI(service)

	paymentV1.RegisterPaymentServiceServer(s, api)

	// Turn on reflection for debugging
	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %s\n", config.AppConfig().PaymentGRPC.Adress())
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
