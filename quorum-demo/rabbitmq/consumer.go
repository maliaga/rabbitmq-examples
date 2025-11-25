package rabbitmq

import (
	"fmt"
	"log"
)

// MessageWithTag represents a message with its delivery tag for manual ack
type MessageWithTag struct {
	Body        string
	DeliveryTag uint64
}

// ConsumeWithManualAck consumes a message without auto-ack
func (r *RabbitMQ) ConsumeWithManualAck() (*MessageWithTag, error) {
	// Get a single message without auto-ack
	msg, ok, err := r.Channel.Get(
		r.QueueName, // queue
		false,       // auto-ack = false (manual acknowledgment)
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume message: %w", err)
	}

	if !ok {
		return nil, fmt.Errorf("no messages available in queue")
	}

	return &MessageWithTag{
		Body:        string(msg.Body),
		DeliveryTag: msg.DeliveryTag,
	}, nil
}

// AckMessage acknowledges a message (confirms successful processing)
func (r *RabbitMQ) AckMessage(deliveryTag uint64) error {
	err := r.Channel.Ack(deliveryTag, false)
	if err != nil {
		return fmt.Errorf("failed to ack message: %w", err)
	}
	log.Printf("✓ Message acknowledged (delivery tag: %d)", deliveryTag)
	return nil
}

// NackMessage negatively acknowledges a message (rejects it)
func (r *RabbitMQ) NackMessage(deliveryTag uint64, requeue bool) error {
	err := r.Channel.Nack(deliveryTag, false, requeue)
	if err != nil {
		return fmt.Errorf("failed to nack message: %w", err)
	}
	if requeue {
		log.Printf("✗ Message rejected and requeued (delivery tag: %d)", deliveryTag)
	} else {
		log.Printf("✗ Message rejected without requeue (delivery tag: %d)", deliveryTag)
	}
	return nil
}

// ConsumeAndAck consumes a message and immediately acknowledges it
func (r *RabbitMQ) ConsumeAndAck() (string, error) {
	msg, err := r.ConsumeWithManualAck()
	if err != nil {
		return "", err
	}

	if err := r.AckMessage(msg.DeliveryTag); err != nil {
		return "", err
	}

	return msg.Body, nil
}

// ConsumeAndNack consumes a message and rejects it (simulates processing failure)
func (r *RabbitMQ) ConsumeAndNack(requeue bool) (string, error) {
	msg, err := r.ConsumeWithManualAck()
	if err != nil {
		return "", err
	}

	if err := r.NackMessage(msg.DeliveryTag, requeue); err != nil {
		return "", err
	}

	return msg.Body, nil
}
