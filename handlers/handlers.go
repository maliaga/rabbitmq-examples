package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"rabbitmq-service/rabbitmq"
)

type Handler struct {
	RabbitMQ *rabbitmq.RabbitMQ
}

type PublishRequest struct {
	Message string `json:"message"`
}

type PublishResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ConsumeResponse struct {
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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Message == "" {
		http.Error(w, "Message cannot be empty", http.StatusBadRequest)
		return
	}

	// Publish message to RabbitMQ
	if err := h.RabbitMQ.PublishMessage(req.Message); err != nil {
		log.Printf("Error publishing message: %v", err)
		http.Error(w, "Failed to publish message", http.StatusInternalServerError)
		return
	}

	log.Printf("Published message: %s", req.Message)

	response := PublishResponse{
		Status:  "success",
		Message: "Message published successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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
		response := ConsumeResponse{
			Status: "error",
			Error:  err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	log.Printf("Consumed message: %s", message)

	response := ConsumeResponse{
		Status:  "success",
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
