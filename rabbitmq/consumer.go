package rabbitmq

import (
	"fmt"
)

// ConsumeMessage consumes a single message from the queue
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
