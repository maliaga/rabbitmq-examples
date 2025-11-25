package rabbitmq

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// SetupQuorumQueue creates a Quorum Queue with replication
func SetupQuorumQueue(ch *amqp.Channel, queueName string) error {
	// Quorum queue arguments
	args := amqp.Table{
		"x-queue-type": "quorum", // Quorum queue type
		// Optional: specify initial group size (default is 3 if cluster has 3+ nodes)
		// "x-quorum-initial-group-size": 3,
	}

	// Declare the Quorum Queue
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable (always true for quorum queues)
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		args,      // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare quorum queue: %w", err)
	}

	log.Printf("✓ Created Quorum Queue: %s", queueName)
	log.Printf("  - Type: Quorum (replicated)")
	log.Printf("  - Durable: true")
	log.Printf("  - Messages: %d", q.Messages)
	log.Printf("  - Consumers: %d", q.Consumers)

	return nil
}

// GetQueueInfo retrieves information about the queue
func GetQueueInfo(ch *amqp.Channel, queueName string) (amqp.Queue, error) {
	// Passive declare to get queue info without creating it
	q, err := ch.QueueDeclarePassive(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return amqp.Queue{}, fmt.Errorf("failed to get queue info: %w", err)
	}

	return q, nil
}
