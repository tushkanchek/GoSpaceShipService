package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	orderV1 "shared/pkg/openapi/order/v1"
	inventoryV1 "shared/pkg/proto/inventory/v1"
	paymentV1 "shared/pkg/proto/payment/v1"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	httpPort          = "8080"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.Order
}


func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.Order),
	}
}

func (s *OrderStorage) GetOrder(order_uuid string) *orderV1.Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[order_uuid]
	if !ok {
		return nil
	}

	return order
}



func (s *OrderStorage) CreateOrder(order *orderV1.Order) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[order.OrderUUID] = order
}


type OrderHandler struct {
	storage *OrderStorage
	inventoryClient inventoryV1.InventoryServiceClient
	paymentClient	paymentV1.PaymentServiceClient
}

func NewOrderHandler(storage *OrderStorage) (*OrderHandler, error) {
	invConn, err := grpc.NewClient(
		"localhost:50051", 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err!=nil{
		return nil, fmt.Errorf("failed connect to inventory: %v", err)
	}

	payConn, err := grpc.NewClient(
		"localhost:50052", 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err!=nil{
		return nil, fmt.Errorf("failed connect to payment: %v", err)
	}

	
	return &OrderHandler{
		storage: storage,
		inventoryClient: inventoryV1.NewInventoryServiceClient(invConn),
		paymentClient: paymentV1.NewPaymentServiceClient(payConn),
	}, err
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error){
	resp, err := h.inventoryClient.ListParts(
		ctx, 
		&inventoryV1.ListPartsRequest{
			Filter: &inventoryV1.PartsFilter{
				Uuids: req.PartUuids,
			},
		},
	)

	if err!=nil{
		return &orderV1.InternalServerError{
			Code: 500,
			Message: fmt.Sprintf("failed to get parts info: %v", err),
		}, nil
	}
	if resp==nil{
		return &orderV1.NotFoundError{
			Code: 404,
			Message: "parts not found",
		}, nil
	}
	if len(resp.Parts)<len(req.PartUuids){
		return &orderV1.BadRequestError{
			Code: 400,
			Message: "some parts not found",
		}, nil
	}

	var totalPrice float64
	for _, p := range resp.Parts{
		totalPrice += p.Price
	}
	
	order := &orderV1.Order{
		OrderUUID: uuid.NewString(),
		UserUUID: req.GetUserUUID(),
		PartUuids: req.GetPartUuids(),
		TotalPrice: float32(totalPrice),
		Status: orderV1.OrderStatusPENDINGPAYMENT,
	}

	h.storage.CreateOrder(order)
	return &orderV1.CreateOrderResponse{
		OrderUUID: order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}, nil
}


func (h *OrderHandler) GetOrderByUUID(_ context.Context, params orderV1.GetOrderByUUIDParams) (orderV1.GetOrderByUUIDRes, error){
	order := h.storage.GetOrder(params.OrderUUID)
	if order==nil{
		return &orderV1.NotFoundError{
			Code: 404,
			Message: "Order for uuid '" + params.OrderUUID + "' not found",
		}, nil
	}

	return &orderV1.GetOrderResponse{
		Order: *order,
	}, nil
}

func (h *OrderHandler) OrderPay(ctx context.Context, req *orderV1.PayOrderRequest,params orderV1.OrderPayParams) (orderV1.OrderPayRes, error){
	order := h.storage.GetOrder(params.OrderUUID)
	if order==nil{
		return &orderV1.NotFoundError{
			Code: 404,
			Message: "order with uuid '" + params.OrderUUID + "' not found",
		}, nil
	}
	transaction_uuid, err := h.paymentClient.PayOrder(
		ctx, 
		&paymentV1.PayOrderRequest{
			OrderUuid: order.OrderUUID,
			UserUuid: order.UserUUID,
			PaymentMethod: string(req.PaymentMethod),
		},
	)
	if err!=nil{
		return nil, fmt.Errorf("payment method error: %v", err)
	}

	order.PaymentMethod = orderV1.NewOptPaymentMethod(req.PaymentMethod)
	order.Status = orderV1.OrderStatusPAID
	order.TransactionUUID = orderV1.NewOptString(transaction_uuid.GetTransactionUuid())
	
	return &orderV1.PayOrderResponse{
		TransactionUUID: transaction_uuid.TransactionUuid,
	}, nil
}



func (h *OrderHandler) OrderCancel(_ context.Context, params orderV1.OrderCancelParams) (orderV1.OrderCancelRes, error){
	order := h.storage.GetOrder(params.OrderUUID)
	if order==nil{
		return &orderV1.NotFoundError{
			Code: 404,
			Message: "order not found",
		}, nil
	}
	if order.Status == orderV1.OrderStatusPAID{
		return &orderV1.ConflictError{
			Code: 409,
			Message: "can't cancel paid order",
		}, nil
	}

	order.Status = orderV1.OrderStatusCANCELLED
	return &orderV1.OrderCancelNoContent{}, nil
}


func (h *OrderHandler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode{
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code: orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}



func main() {
	storage := NewOrderStorage()

	orderHandler, err := NewOrderHandler(storage)
	if err!=nil{
		log.Fatalf("Ошибка при создании OrderHandler: %v", err)
	}

	orderServer, err := orderV1.NewServer(orderHandler)
	if err!=nil{
		log.Fatalf("Ошибка при создании сервера OpenAPI: %v", err)
	}
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr: 	net.JoinHostPort("localhost", httpPort),
		Handler: r,
		ReadHeaderTimeout: readHeaderTimeout,
	}


	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")

}
