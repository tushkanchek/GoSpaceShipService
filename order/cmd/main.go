package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	orderV1API "order/internal/api/order/v1"
	clientInventory "order/internal/client/grpc/inventory/v1"
	clientPayment "order/internal/client/grpc/payment/v1"
	"order/internal/migrator"
	repoOrder "order/internal/repository/order"
	serviceOrder "order/internal/service/order"
	orderV1 "shared/pkg/openapi/order/v1"
	inventoryV1 "shared/pkg/proto/inventory/v1"
	paymentV1 "shared/pkg/proto/payment/v1"
)

const (
	httpPort          = "8080"
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
	inventoryAdress   = "50051"
	paymentAdress     = "50052"

	envPathDefault = ".env"
	DB_URI         = "DB_URI"
	MIGRATIONS_DIR = "MIGRATIONS_DIR"
)

func main() {
	ctx := context.Background()

	err := godotenv.Load(envPathDefault)
	if err != nil {
		log.Printf("failed to load .env file: %v\n", err)
		return
	}

	dbURI := os.Getenv(DB_URI)

	con, err := pgx.Connect(ctx, dbURI)
	if err != nil {
		log.Printf("failed connect to db: %v\n", err)
		return
	}
	defer func() {
		cerr := con.Close(ctx)
		if cerr != nil {
			log.Printf("failed to close connection: %v\n", cerr)
		}
	}()

	err = con.Ping(ctx)
	if err != nil {
		log.Printf("data base is unavailable")
		return
	}

	migrationDir := os.Getenv(MIGRATIONS_DIR)
	migratorRunner := migrator.NewMigrator(stdlib.OpenDB(*con.Config().Copy()), migrationDir)

	err = migratorRunner.Up()
	if err != nil {
		log.Printf("Migration error: %v\n", err)
	}

	invConn, err := grpc.NewClient(
		"localhost:"+inventoryAdress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("Failed connect to inventory: %v\n", err)
		return
	}

	invClient := inventoryV1.NewInventoryServiceClient(invConn)

	payConn, err := grpc.NewClient(
		"localhost:"+paymentAdress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("Failed connect to payment: %v\n", err)
		return
	}

	payClient := paymentV1.NewPaymentServiceClient(payConn)

	gprcInventory := clientInventory.NewClient(invClient)

	gprcPayment := clientPayment.NewClient(shutdownTimeout, payClient)

	log.Printf("Succesfully created grpc clients")

	repo := repoOrder.NewOrderRepository(con) // Подумать про реализацию Pool

	service := serviceOrder.NewService(repo, gprcInventory, gprcPayment)

	api := orderV1API.NewAPI(service)

	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		log.Printf("Ошибка при создании сервера OpenAPI: %v", err)
		return
	}
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
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
