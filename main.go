package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rabbitmq-service/handlers"
	"rabbitmq-service/rabbitmq"
	"syscall"
)

func main() {
	// Get configuration from environment variables or use defaults
	rabbitmqURL := getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	queueName := getEnv("RABBITMQ_QUEUE_NAME", "messages")
	httpPort := getEnv("HTTP_PORT", "8080")

	// Initialize RabbitMQ connection
	rmq, err := rabbitmq.NewRabbitMQ(rabbitmqURL, queueName)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer rmq.Close()

	// Create handler with RabbitMQ instance
	handler := &handlers.Handler{
		RabbitMQ: rmq,
	}

	// Setup HTTP routes
	http.HandleFunc("/publish", handler.PublishHandler)
	http.HandleFunc("/consume", handler.ConsumeHandler)
	http.HandleFunc("/health", healthHandler)

	// Start HTTP server in a goroutine
	go func() {
		addr := fmt.Sprintf(":%s", httpPort)
		log.Printf("Starting HTTP server on %s", addr)
		log.Printf("Endpoints:")
		log.Printf("  POST http://localhost:%s/publish - Publish a message", httpPort)
		log.Printf("  GET  http://localhost:%s/consume - Consume a message", httpPort)
		log.Printf("  GET  http://localhost:%s/health  - Health check", httpPort)
		log.Printf("\nRabbitMQ Management UI: http://localhost:15672 (guest/guest)")
		
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
