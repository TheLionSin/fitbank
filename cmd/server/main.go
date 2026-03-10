package main

import (
	"context"
	"fitbank/internal/handler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	activityHandler := handler.NewActivityHandler()

	mux := http.NewServeMux()

	// Наш тестовый эндпоинт
	mux.HandleFunc("POST /activities", activityHandler.Create)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	// 2. Настраиваем параметры сервера
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
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
