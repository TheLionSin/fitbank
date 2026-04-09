package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// Наши внутренние пакеты
	"fitbank/nutrition-service/internal/app"
	"fitbank/nutrition-service/internal/client"
	"fitbank/nutrition-service/internal/repository"
	"fitbank/nutrition-service/pkg/api" // Сгенерированный код этого сервиса

	// Сторонние библиотеки
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // Нужно для goose
	"github.com/pressly/goose/v3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// 1. Инициализируем логгер (в банках обычно используют структурированный логгер slog)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// 2. Читаем конфигурацию из переменных окружения
	// В реальном проде здесь будет библиотека типа viper или cleanenv
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		// Для локальной разработки, если переменная не задана
		dbURL = "postgres://user:pass@localhost:5433/nutrition_db?sslmode=disable"
	}

	activityAddr := os.Getenv("ACTIVITY_SERVICE_ADDR")
	if activityAddr == "" {
		// Адрес gRPC сервера первого микросервиса
		activityAddr = "localhost:50051"
	}

	// 3. Создаем контекст, который отменится при нажатии Ctrl+C
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 4. Запускаем миграции (Goose)
	runMigrations(dbURL)

	// 5. Подключаемся к базе данных через пул соединений
	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer dbPool.Close()

	// 6. Инициализируем gRPC-клиент для связи с Activity Service
	actClient, err := client.NewActivityClient(activityAddr)
	if err != nil {
		log.Fatalf("failed to create activity client: %v", err)
	}
	// Важно: в реальном клиенте нужно не забыть закрыть соединение при выходе,
	// но для простоты пока оставим так.

	// 7. Собираем слои приложения (Dependency Injection)
	repo := repository.NewPostgresMealRepository(dbPool)
	nutritionService := app.NewService(repo, actClient)

	// Создаем gRPC сервер нашего сервиса (реализуем позже)
	grpcServer := app.NewNutritionGRPCServer(nutritionService)

	// 8. Запуск серверов через errgroup
	g, ctx := errgroup.WithContext(ctx)

	// --- ГОРУТИНА: HTTP Server (Метрики) ---
	g.Go(func() error {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())

		srv := &http.Server{
			Addr:    ":8081", // Порт 8081, так как 8080 занят первым сервисом
			Handler: mux,
		}

		slog.Info("Metrics server starting", "addr", srv.Addr)

		go func() {
			<-ctx.Done()
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			srv.Shutdown(shutdownCtx)
		}()

		return srv.ListenAndServe()
	})

	// --- ГОРУТИНА: gRPC Server (Основная логика) ---
	g.Go(func() error {
		lis, err := net.Listen("tcp", ":50052") // Порт 50052
		if err != nil {
			return err
		}

		s := grpc.NewServer()
		// РЕГИСТРИРУЕМ НАШ СЕРВИС:
		nutritionpb.RegisterNutritionServiceServer(s, grpcServer)
		reflection.Register(s) // Нужно для отладки через Postman/Evans

		slog.Info("gRPC server starting", "addr", lis.Addr().String())

		go func() {
			<-ctx.Done()
			slog.Info("gRPC server shutting down")
			s.GracefulStop()
		}()

		return s.Serve(lis)
	})

	// Ждем завершения всех горутин
	if err := g.Wait(); err != nil {
		slog.Error("service finished with error", "error", err)
	}
}

// Вспомогательная функция для миграций
func runMigrations(dbURL string) {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("failed to open db for migrations: %v", err)
	}
	defer db.Close()

	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}
	slog.Info("migrations completed successfully")
}
