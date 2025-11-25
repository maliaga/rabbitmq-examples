package rabbitmq

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	QueueName  string
	DLQName    string
}

// NewRabbitMQWithDLX creates a new RabbitMQ connection with Dead Letter Exchange support
func NewRabbitMQWithDLX(url, queueName string) (*RabbitMQ, error) {
	// Connect to RabbitMQ
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Setup Dead Letter Exchange and Dead Letter Queue
	if err := SetupDLX(ch); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to setup DLX: %w", err)
	}

	// Setup main queue with DLX configuration
	if err := SetupMainQueueWithDLX(ch, queueName); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to setup main queue with DLX: %w", err)
	}

	log.Printf("Connected to RabbitMQ with DLX enabled")
	log.Printf("Main Queue: %s", queueName)
	log.Printf("Dead Letter Queue: %s", DLQName)

	return &RabbitMQ{
		Connection: conn,
		Channel:    ch,
		QueueName:  queueName,
		DLQName:    DLQName,
	}, nil
}

// Close closes the channel and connection
func (r *RabbitMQ) Close() {
	if r.Channel != nil {
		r.Channel.Close()
	}
	if r.Connection != nil {
		r.Connection.Close()
	}
	log.Println("RabbitMQ connection closed")
}
