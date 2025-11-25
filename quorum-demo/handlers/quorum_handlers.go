package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"rabbitmq-quorum-demo/rabbitmq"
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
	Data    any    `json:"data,omitempty"`
}

// PublishHandler handles POST requests to publish messages with confirmation
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

	// Publish message with confirmation
	if err := h.RabbitMQ.PublishWithConfirmation(req.Message); err != nil {
		log.Printf("Error publishing message: %v", err)
		respondWithError(w, "Failed to publish message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Published and confirmed: %s", req.Message)
	respondWithSuccess(w, "Message published and confirmed by broker")
}

// ConsumeHandler handles GET requests to consume messages with ACK
func (h *Handler) ConsumeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Consume message and acknowledge
	message, err := h.RabbitMQ.ConsumeAndAck()
	if err != nil {
		log.Printf("Error consuming message: %v", err)
		respondWithError(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Printf("Consumed and acknowledged: %s", message)
	respondWithMessage(w, message)
}

// ConsumeWithFailureHandler handles POST requests to consume and NACK messages
func (h *Handler) ConsumeWithFailureHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Consume message and reject it (simulate processing failure)
	message, err := h.RabbitMQ.ConsumeAndNack(true) // requeue = true
	if err != nil {
		log.Printf("Error consuming message: %v", err)
		respondWithError(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Printf("Consumed and rejected (requeued): %s", message)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Status:  "success",
		Message: "Message consumed and rejected (requeued): " + message,
	})
}

// StatsHandler handles GET requests to show queue statistics
func (h *Handler) StatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get queue info
	queueInfo, err := rabbitmq.GetQueueInfo(h.RabbitMQ.Channel, h.RabbitMQ.QueueName)
	if err != nil {
		log.Printf("Error getting queue info: %v", err)
		respondWithError(w, "Failed to get queue stats", http.StatusInternalServerError)
		return
	}

	stats := map[string]interface{}{
		"queue_name": h.RabbitMQ.QueueName,
		"queue_type": "quorum",
		"messages":   queueInfo.Messages,
		"consumers":  queueInfo.Consumers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Status: "success",
		Data:   stats,
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
