package msg_broker

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func New() *amqp.Connection {
	conn, err := amqp.Dial("amqp://guest:guest@172.53.1.10:5672/")
	if err != nil {
		log.Fatalf("Failed to establish RabbitMQ connection: %v", err.Error())
	}
	return conn
}
