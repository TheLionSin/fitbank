package main

import (
	"context"
	"database/sql"
	"fitbank/activity-service/internal/config"
	"fitbank/activity-service/internal/handler"
	"fitbank/activity-service/internal/middleware"
	"fitbank/activity-service/internal/repository"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
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

	migrationDB, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("failed to open db for migrations: %v", err)
	}

	log.Println("Running migrations...")
	if err := goose.Up(migrationDB, "migrations"); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}
	migrationDB.Close() // Закрываем его, он больше не нужен
	log.Println("Migrations finished!")

	repo := repository.NewPostgresRepository(dbPool)
	activityHandler := handler.NewActivityHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /activities", activityHandler.Create)
	mux.HandleFunc("GET /activities", activityHandler.List)
	mux.HandleFunc("GET /activities/{id}", activityHandler.Get)
	mux.HandleFunc("PUT /activities/{id}", activityHandler.Update)
	mux.HandleFunc("DELETE /activities/{id}", activityHandler.Delete)

	// Оборачиваем mux в middleware (цепочка)
	// Сначала назначаем ID, потом логируем
	var handler http.Handler = mux
	handler = middleware.Logging(handler)
	handler = middleware.RequestID(handler)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	// 2. Настраиваем параметры сервера
	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,  //сколько времени мы даем клиенту на отправку запроса.
		WriteTimeout: 10 * time.Second, //сколько времени мы готовы ждать, пока клиент скачивает наш ответ.
		IdleTimeout:  120 * time.Second,
	}

	// 3. Канал для перехвата сигналов прерывания (Ctrl+C, kill)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Запускаем сервер в отдельной горутине, чтобы он не блокировал основной поток
	go func() {
		log.Printf("Activity Service starting on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
		}
	}()

	// Ждем сигнала остановки
	<-stop
	log.Println("Shutting down server...")

	// 4. Graceful Shutdown: даем серверу 5 секунд на завершение дел
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}
	log.Println("Server gracefully stopped")
}
