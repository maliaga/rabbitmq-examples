package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rabbitmq-quorum-demo/handlers"
	"rabbitmq-quorum-demo/rabbitmq"
	"syscall"
)

func main() {
	// Get configuration from environment variables or use defaults
	rabbitmqURL := getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	queueName := getEnv("RABBITMQ_QUEUE_NAME", "orders-quorum")
	httpPort := getEnv("HTTP_PORT", "8082")

	// Initialize RabbitMQ connection with Quorum Queue support
	rmq, err := rabbitmq.NewRabbitMQWithQuorum(rabbitmqURL, queueName)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ with Quorum Queue: %v", err)
	}
	defer rmq.Close()

	// Create handler with RabbitMQ instance
	handler := &handlers.Handler{
		RabbitMQ: rmq,
	}

	// Setup HTTP routes
	http.HandleFunc("/publish", handler.PublishHandler)
	http.HandleFunc("/consume", handler.ConsumeHandler)
	http.HandleFunc("/consume/fail", handler.ConsumeWithFailureHandler)
	http.HandleFunc("/stats", handler.StatsHandler)
	http.HandleFunc("/health", healthHandler)

	// Start HTTP server in a goroutine
	go func() {
		addr := fmt.Sprintf(":%s", httpPort)
		log.Printf("========================================")
		log.Printf("RabbitMQ Quorum Queue Demo Service")
		log.Printf("========================================")
		log.Printf("HTTP Server: http://localhost:%s", httpPort)
		log.Printf("")
		log.Printf("Endpoints:")
		log.Printf("  POST http://localhost:%s/publish       - Publish with confirmation", httpPort)
		log.Printf("  GET  http://localhost:%s/consume       - Consume with ACK", httpPort)
		log.Printf("  POST http://localhost:%s/consume/fail  - Consume with NACK (requeue)", httpPort)
		log.Printf("  GET  http://localhost:%s/stats         - Queue statistics", httpPort)
		log.Printf("  GET  http://localhost:%s/health        - Health check", httpPort)
		log.Printf("")
		log.Printf("RabbitMQ Cluster Management UIs:")
		log.Printf("  Node 1: http://localhost:15672 (guest/guest)")
		log.Printf("  Node 2: http://localhost:15673 (guest/guest)")
		log.Printf("  Node 3: http://localhost:15674 (guest/guest)")
		log.Printf("========================================")

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
	w.Write([]byte(`{"status":"healthy","queue_type":"quorum"}`))
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
