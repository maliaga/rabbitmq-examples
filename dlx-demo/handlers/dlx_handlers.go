package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"rabbitmq-dlx-demo/rabbitmq"
)

type Handler struct {
	RabbitMQ *rabbitmq.RabbitMQ
}

type PublishRequest struct {
	Message string `json:"message"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// PublishHandler handles POST requests to publish messages
func (h *Handler) PublishHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req PublishRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Message == "" {
		respondWithError(w, "Message cannot be empty", http.StatusBadRequest)
		return
	}

	// Publish message to RabbitMQ
	if err := h.RabbitMQ.PublishMessage(req.Message); err != nil {
		log.Printf("Error publishing message: %v", err)
		respondWithError(w, "Failed to publish message", http.StatusInternalServerError)
		return
	}

	log.Printf("Published message: %s", req.Message)
	respondWithSuccess(w, "Message published successfully")
}

// ConsumeHandler handles GET requests to consume messages
func (h *Handler) ConsumeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Consume message from RabbitMQ
	message, err := h.RabbitMQ.ConsumeMessage()
	if err != nil {
		log.Printf("Error consuming message: %v", err)
		respondWithError(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Printf("Consumed message: %s", message)
	respondWithMessage(w, message)
}

// RejectMessageHandler handles POST requests to reject messages (simulate failure)
func (h *Handler) RejectMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Consume and reject a message (sends to DLX)
	message, err := h.RabbitMQ.RejectMessage()
	if err != nil {
		log.Printf("Error rejecting message: %v", err)
		respondWithError(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Printf("Rejected message (sent to DLX): %s", message)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Status:  "success",
		Message: "Message rejected and sent to DLX: " + message,
	})
}

// ConsumeDLQHandler handles GET requests to consume from Dead Letter Queue
func (h *Handler) ConsumeDLQHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Consume message from DLQ
	message, err := h.RabbitMQ.ConsumeFromDLQ()
	if err != nil {
		log.Printf("Error consuming from DLQ: %v", err)
		respondWithError(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Printf("Consumed from DLQ: %s", message)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Status:  "success",
		Message: "Message from DLQ: " + message,
	})
}

// Helper functions
func respondWithError(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(Response{
		Status: "error",
		Error:  error,
	})
}

func respondWithSuccess(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Status:  "success",
		Message: message,
	})
}

func respondWithMessage(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Status:  "success",
		Message: message,
	})
}
