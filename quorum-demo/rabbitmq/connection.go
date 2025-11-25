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
}

// NewRabbitMQWithQuorum creates a new RabbitMQ connection with Quorum Queue support
func NewRabbitMQWithQuorum(url, queueName string) (*RabbitMQ, error) {
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

	// Enable publisher confirmations
	if err := ch.Confirm(false); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to enable publisher confirmations: %w", err)
	}
	log.Println("✓ Publisher confirmations enabled")

	// Setup Quorum Queue
	if err := SetupQuorumQueue(ch, queueName); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to setup quorum queue: %w", err)
	}

	log.Printf("✓ Connected to RabbitMQ with Quorum Queue support")
	log.Printf("  Queue: %s (type: quorum)", queueName)

	return &RabbitMQ{
		Connection: conn,
		Channel:    ch,
		QueueName:  queueName,
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
