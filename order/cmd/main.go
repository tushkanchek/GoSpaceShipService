package main

import (
	"context"
	"errors"
	"fmt"
	"log"
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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	orderV1API "order/internal/api/order/v1"
	clientInventory "order/internal/client/grpc/inventory/v1"
	clientPayment "order/internal/client/grpc/payment/v1"
	"order/internal/config"
	"order/internal/migrator"
	repoOrder "order/internal/repository/order"
	serviceOrder "order/internal/service/order"
	orderV1 "shared/pkg/openapi/order/v1"
	inventoryV1 "shared/pkg/proto/inventory/v1"
	paymentV1 "shared/pkg/proto/payment/v1"
)

const (
	configPath = "./deploy/compose/order/.env"
)

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create connection to postgres
	con, err := pgx.Connect(ctx, config.AppConfig().Postgres.URI())
	if err != nil {
		log.Printf("failed connect to db: %v\n", err)
		return
	}
	defer func() {
		cerr := con.Close(context.Background())
		if cerr != nil {
			log.Printf("failed to close connection: %v\n", cerr)
		}
	}()

	// Check the connection to Postgresql
	err = con.Ping(ctx)
	if err != nil {
		log.Printf("data base is unavailable: %v\n", err)
		return
	}
	log.Println("✅ Connected to postgres")

	// Activate migrations
	migratorRunner := migrator.NewMigrator(stdlib.OpenDB(*con.Config().Copy()), config.AppConfig().Postgres.MigrationDir())

	err = migratorRunner.Up()
	if err != nil {
		log.Printf("Migration error: %v\n", err)
	}

	// Create Inventory GRPC service client
	invConn, err := grpc.NewClient(
		config.AppConfig().InventoryGRPC.Adress(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("Failed connect to inventory: %v\n", err)
		return
	}
	invClient := inventoryV1.NewInventoryServiceClient(invConn)

	// Create Payment GRPC service client
	payConn, err := grpc.NewClient(
		config.AppConfig().PaymentGRPC.Adress(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("Failed connect to payment: %v\n", err)
		return
	}
	payClient := paymentV1.NewPaymentServiceClient(payConn)

	// Create GRPC Clients
	gprcInventory := clientInventory.NewClient(invClient)

	gprcPayment := clientPayment.NewClient(payClient)
	log.Printf("✅ Succesfully created grpc clients")

	// Register OpenAPI server
	repo := repoOrder.NewOrderRepository(con) // Подумать про реализацию Pool

	service := serviceOrder.NewService(repo, gprcInventory, gprcPayment)

	api := orderV1API.NewAPI(service)

	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		log.Printf("Ошибка при создании сервера OpenAPI: %v", err)
		return
	}

	// Create and configure router
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              config.AppConfig().OrderHTTP.Adress(),
		Handler:           r,
		ReadHeaderTimeout: config.AppConfig().OrderHTTP.ReadTimeout(),
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("🚀 HTTP-сервер запущен на %s\n", config.AppConfig().OrderHTTP.Adress())
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
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
