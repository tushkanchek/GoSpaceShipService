package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"order/internal/config"
	"order/internal/migrator"
	"platform/pkg/closer"
	"platform/pkg/logger"
	orderV1 "shared/pkg/openapi/order/v1"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type App struct {
	diContainer *diContainer
	httpServer  *http.Server
	router      *chi.Mux
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	// HTTP server
	g.Go(func() error {
		return a.runHTTPServer(ctx)
	})

	// Kafka consumer
	g.Go(func() error {
		return a.runConsumer(ctx)
	})

	return g.Wait()
}



func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initMigrations,
		a.initRouter,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initMigrations(ctx context.Context) error {
	// Use separate connection for migrations
	migConn, err := pgx.Connect(ctx, config.AppConfig().Postgres.URI())
	if err != nil {
		return fmt.Errorf("failed to connect for migrations: %w", err)
	}

	closer.AddNamed("Migration connection", func(ctx context.Context) error {
		return migConn.Close(ctx)
	})

	migratorRunner := migrator.NewMigrator(
		stdlib.OpenDB(*migConn.Config().Copy()),
		config.AppConfig().Postgres.MigrationDir(),
	)

	err = migratorRunner.Up()
	if err != nil {
		return fmt.Errorf("migration error: %w", err)
	}

	return nil
}

func (a *App) initRouter(ctx context.Context) error {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Create OpenAPI server
	orderServer, err := orderV1.NewServer(a.diContainer.OrderV1API(ctx))
	if err != nil {
		return fmt.Errorf("failed to create OpenAPI server: %w", err)
	}

	r.Mount("/", orderServer)

	a.router = r
	return nil
}

func (a *App) initHTTPServer(_ context.Context) error {
	a.httpServer = &http.Server{
		Addr:              config.AppConfig().OrderHTTP.Adress(),
		Handler:           a.router,
		ReadHeaderTimeout: config.AppConfig().OrderHTTP.ReadTimeout(),
	}

	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		return a.httpServer.Shutdown(ctx)
	})

	return nil
}



func (a *App) runHTTPServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 HTTP OrderService server listening on %s", config.AppConfig().OrderHTTP.Adress()))
	errCh := make(chan error, 1)

	go func() {
		err := a.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		logger.Error(ctx, "❌ Context done error %v", zap.Error(ctx.Err()))
		return nil
	case err := <-errCh:
		logger.Error(ctx, "❌ HTTP server error: %v", zap.Error(err))
		return err
	}
}


func (a *App) runConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 Kafka consumer started")

	errCh := make(chan error, 1)

	go func() {
		err := a.diContainer.ConsumerService(ctx).RunConsumer(ctx)
		if err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info(ctx, "🛑 Kafka consumer shutdown signal received")
		return nil
	case err := <-errCh:
		logger.Error(ctx, "❌ Kafka consumer error", zap.Error(err))
		return err
	}
}