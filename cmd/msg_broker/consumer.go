package msg_broker

import (
	"context"
	"encoding/json"
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
		"new_post_queue", // name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
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

	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	messages, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto ack
		false,      // exclusive
		false,      // no local
		false,      // no wait
		nil,        // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {

		var msgPayload contracts.MessagePayload

		//consume messages
		for d := range messages {
			log.Printf(" [x] %s", d.Body)

			//Decode struct payload
			err := json.Unmarshal(d.Body, &msgPayload)
			if err != nil {
				log.Println(err.Error())
			}

			_, err = c.notificationService.Broadcast(context.TODO(), "newPost", &models.Message{
				Title:    msgPayload.Title,
				Body:     msgPayload.Body,
				ImageUrl: msgPayload.ImageUrl,
			})
			if err != nil {
				log.Println(err.Error())
			}

			err = d.Ack(false)
			if err != nil {
				log.Println(err.Error())
			}

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
