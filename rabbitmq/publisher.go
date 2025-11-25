package rabbitmq

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// PublishMessage publishes a message to the queue
func (r *RabbitMQ) PublishMessage(message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.Channel.PublishWithContext(
		ctx,
		"",           // exchange
		r.QueueName,  // routing key (queue name)
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
			DeliveryMode: amqp.Persistent, // make message persistent
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}
