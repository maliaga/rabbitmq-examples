package rabbitmq

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	DLXExchangeName = "dlx.exchange"
	DLQName         = "messages-dlx.dlq"
	DLXRoutingKey   = "dlx.routing.key"
)

// SetupDLX creates the Dead Letter Exchange and Dead Letter Queue
func SetupDLX(ch *amqp.Channel) error {
	// Declare the Dead Letter Exchange
	err := ch.ExchangeDeclare(
		DLXExchangeName, // name
		"direct",        // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare DLX exchange: %w", err)
	}
	log.Printf("Created Dead Letter Exchange: %s", DLXExchangeName)

	// Declare the Dead Letter Queue
	_, err = ch.QueueDeclare(
		DLQName, // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare DLQ: %w", err)
	}
	log.Printf("Created Dead Letter Queue: %s", DLQName)

	// Bind the DLQ to the DLX
	err = ch.QueueBind(
		DLQName,         // queue name
		DLXRoutingKey,   // routing key
		DLXExchangeName, // exchange
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind DLQ to DLX: %w", err)
	}
	log.Printf("Bound DLQ to DLX with routing key: %s", DLXRoutingKey)

	return nil
}

// SetupMainQueueWithDLX creates the main queue with DLX configuration
func SetupMainQueueWithDLX(ch *amqp.Channel, queueName string) error {
	// Queue arguments to enable DLX
	args := amqp.Table{
		"x-dead-letter-exchange":    DLXExchangeName,
		"x-dead-letter-routing-key": DLXRoutingKey,
		// Optional: set message TTL (time to live) in milliseconds
		// "x-message-ttl": 60000, // 60 seconds
	}

	// Declare the main queue with DLX arguments
	_, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		args,      // arguments with DLX configuration
	)
	if err != nil {
		return fmt.Errorf("failed to declare main queue with DLX: %w", err)
	}
	log.Printf("Created main queue with DLX: %s", queueName)

	return nil
}
