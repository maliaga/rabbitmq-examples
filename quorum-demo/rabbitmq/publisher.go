package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// PublishWithConfirmation publishes a message and waits for broker confirmation
func (r *RabbitMQ) PublishWithConfirmation(message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a channel to receive confirmations
	confirms := r.Channel.NotifyPublish(make(chan amqp.Confirmation, 1))

	// Publish the message
	err := r.Channel.PublishWithContext(
		ctx,
		"",          // exchange
		r.QueueName, // routing key (queue name)
		true,        // mandatory - return message if not routable
		false,       // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(message),
			DeliveryMode: amqp.Persistent, // persistent (required for quorum queues)
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	// Wait for confirmation
	select {
	case confirm := <-confirms:
		if confirm.Ack {
			log.Printf("✓ Message confirmed by broker: %s", message)
			return nil
		}
		return fmt.Errorf("message not confirmed (nack received)")
	case <-ctx.Done():
		return fmt.Errorf("timeout waiting for confirmation")
	}
}

// PublishBatch publishes multiple messages with confirmations
func (r *RabbitMQ) PublishBatch(messages []string) (int, error) {
	successCount := 0
	for i, msg := range messages {
		if err := r.PublishWithConfirmation(msg); err != nil {
			log.Printf("✗ Failed to publish message %d: %v", i+1, err)
			return successCount, err
		}
		successCount++
	}
	return successCount, nil
}
