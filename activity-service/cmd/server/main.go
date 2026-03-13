package main

import (
	"context"
	"database/sql"
	"fitbank/activity-service/internal/app"
	"fitbank/activity-service/internal/config"
	"fitbank/activity-service/internal/handler"
	"fitbank/activity-service/internal/middleware"
	"fitbank/activity-service/internal/repository"
	activitypb "fitbank/activity-service/pkg/api"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func main() {

	cfg := config.New()

	// Настраиваем логгер: в проде JSON, в деве — обычный текст
	var logger *slog.Logger
	if cfg.Environment == "prod" {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
	slog.SetDefault(logger)

	slog.Info("Starting activity-service", "port", cfg.AppPort, "env", cfg.Environment)

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/fitbank?sslmode=disable" // дефолт для локального запуска
	}

	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	var migrationDB *sql.DB

	for i := 0; i < 5; i++ {
		migrationDB, err = sql.Open("pgx", dbURL)
		if err == nil {
			err = migrationDB.Ping()
		}

		if err == nil {
			break
		}

		if migrationDB != nil {
			migrationDB.Close()
		}

		slog.Warn("DB not ready, retrying...", "attempt", i+1, "error", err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("failed to connect to db after retries: %v", err)
	}

	log.Println("Running migrations...")
	if err := goose.Up(migrationDB, "migrations"); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}
	migrationDB.Close() // Закрываем его, он больше не нужен
	slog.Info("Migrations finished!")

	repo := repository.NewPostgresRepository(dbPool)
	actService := app.NewService(repo)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		mux := http.NewServeMux()
		h := handler.NewActivityHandler(actService)

		mux.HandleFunc("POST /activities", h.Create)
		mux.HandleFunc("GET /activities", h.List)
		mux.HandleFunc("GET /activities/{id}", h.Get)
		mux.HandleFunc("PUT /activities/{id}", h.Update)
		mux.HandleFunc("DELETE /activities/{id}", h.Delete)
		mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"status":"ok"}`))
		})

		// Оборачиваем в Middleware
		var finalHandler http.Handler = mux
		finalHandler = middleware.Logging(finalHandler)
		finalHandler = middleware.RequestID(finalHandler)

		server := &http.Server{
			Addr:    ":8080",
			Handler: finalHandler,
		}

		// Запускаем горутину для мягкой остановки именно этого сервера
		go func() {
			<-ctx.Done() // Ждем сигнала отмены
			slog.Info("Shutting down HTTP server...")
			server.Shutdown(context.Background())
		}()

		slog.Info("HTTP server starting", "port", "8080")
		return server.ListenAndServe()
	})

	// 3. gRPC Server
	g.Go(func() error {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			return err
		}

		s := grpc.NewServer()
		activitypb.RegisterActivityServiceServer(s, app.NewGRPCServer(actService))

		go func() {
			<-ctx.Done()
			slog.Info("Shutting down gRPC server...")
			s.GracefulStop()
		}()

		slog.Info("gRPC server starting", "port", "50051")
		return s.Serve(lis)
	})

	// Ждем завершения
	if err := g.Wait(); err != nil && err != http.ErrServerClosed {
		slog.Error("Service stopped", "error", err)
	}
	slog.Info("Service gracefully stopped")

}
