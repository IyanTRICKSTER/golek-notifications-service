package msg_broker

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"golek_notifications_service/pkg/contracts"
	"golek_notifications_service/pkg/models"
	"log"
)

type MQConsumer struct {
	notificationService contracts.NotificationService
	amqpConnection      *amqp.Connection
}

func (c *MQConsumer) Listen() {

	//Close connection
	defer c.amqpConnection.Close()

	//open channel
	channel, err := c.amqpConnection.Channel()
	failOnError(err, "Failed to open a channel!")
	defer channel.Close()

	//declare exchange
	err = channel.ExchangeDeclare(
		"posts_exchange", // name
		"topic",          // type
		true,             // durable
		false,            // auto-deleted
		false,            // internal
		false,            // no-wait
		nil,
	)
	failOnError(err, "Failed to declare an exchange!")

	queue, err := channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue!")

	err = channel.QueueBind(
		queue.Name,
		"new_post_route",
		"posts_exchange",
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	messages, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto ack
		false,      // exclusive
		false,      // no local
		false,      // no wait
		nil,        // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		err := func() error {
			for d := range messages {
				log.Printf(" [x] %s", d.Body)
				_, err := c.notificationService.Broadcast(context.TODO(), "newPost", &models.Message{
					Title:    string(d.Body),
					Body:     "",
					ImageUrl: "",
				})
				return err
			}
			return nil
		}()
		if err != nil {
			panic(err)
		}
	}()

	log.Printf(" [*] Waiting for message. To exit press CTRL+C")
	<-forever
}

func NewConsumer(notifyService *contracts.NotificationService, amqpConn *amqp.Connection) *MQConsumer {
	return &MQConsumer{
		notificationService: *notifyService,
		amqpConnection:      amqpConn,
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
