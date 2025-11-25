package rabbitmq

import (
	"fmt"
	"log"
)

// ConsumeMessage consumes a single message from the queue with auto-ack
func (r *RabbitMQ) ConsumeMessage() (string, error) {
	// Get a single message
	msg, ok, err := r.Channel.Get(
		r.QueueName, // queue
		true,        // auto-ack
	)
	if err != nil {
		return "", fmt.Errorf("failed to consume message: %w", err)
	}

	if !ok {
		return "", fmt.Errorf("no messages available in queue")
	}

	return string(msg.Body), nil
}

// ConsumeMessageManual consumes a message without auto-ack (for manual ack/nack)
func (r *RabbitMQ) ConsumeMessageManual() (string, uint64, error) {
	// Get a single message without auto-ack
	msg, ok, err := r.Channel.Get(
		r.QueueName, // queue
		false,       // auto-ack = false (manual acknowledgment)
	)
	if err != nil {
		return "", 0, fmt.Errorf("failed to consume message: %w", err)
	}

	if !ok {
		return "", 0, fmt.Errorf("no messages available in queue")
	}

	return string(msg.Body), msg.DeliveryTag, nil
}

// RejectMessage consumes a message and rejects it (sends to DLX)
func (r *RabbitMQ) RejectMessage() (string, error) {
	// Get a message without auto-ack
	message, deliveryTag, err := r.ConsumeMessageManual()
	if err != nil {
		return "", err
	}

	// Reject the message (nack with requeue=false sends it to DLX)
	err = r.Channel.Nack(
		deliveryTag, // delivery tag
		false,       // multiple
		false,       // requeue = false (send to DLX instead of requeuing)
	)
	if err != nil {
		return "", fmt.Errorf("failed to reject message: %w", err)
	}

	log.Printf("Message rejected and sent to DLX: %s", message)
	return message, nil
}

// ConsumeFromDLQ consumes a message from the Dead Letter Queue
func (r *RabbitMQ) ConsumeFromDLQ() (string, error) {
	// Get a single message from DLQ
	msg, ok, err := r.Channel.Get(
		r.DLQName, // dead letter queue
		true,      // auto-ack
	)
	if err != nil {
		return "", fmt.Errorf("failed to consume from DLQ: %w", err)
	}

	if !ok {
		return "", fmt.Errorf("no messages available in DLQ")
	}

	return string(msg.Body), nil
}
