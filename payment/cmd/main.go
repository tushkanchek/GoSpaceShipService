package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	paymentV1 "shared/pkg/proto/payment/v1"
	"syscall"

	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)




const (
	grpcPort = 50052
)


type PaymentService struct{
	paymentV1.UnimplementedPaymentServiceServer
}

//TODO: validate PaymentMethod
func (s *PaymentService) PayOrder(_ context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error){
	transaction_uuid := uuid.NewString()
	log.Printf("Оплата успешно прошла, transaction_uuid: %s\n", transaction_uuid)
	return &paymentV1.PayOrderResponse{
		TransactionUuid: transaction_uuid,
	}, nil
}


func main(){
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err!=nil{
		log.Printf("Failed to listen: %v\n", err)
		return
	}
	defer func(){
		if cerr:=lis.Close();cerr!=nil{
			log.Printf("Failed to close listener: %v\n", cerr)
		}
	}()
	
	s:=grpc.NewServer()

	service := &PaymentService{}

	paymentV1.RegisterPaymentServiceServer(s, service)

	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
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

