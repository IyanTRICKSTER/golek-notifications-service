package main

import (
	"context"
	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"golek_notifications_service/cmd/msg_broker"
	"golek_notifications_service/pkg/repositories"
	"golek_notifications_service/pkg/services"
	"google.golang.org/api/option"

	"log"
	"os"
)

func main() {

	//load .env
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	//load fcm key
	wd, _ := os.Getwd()
	opt := option.WithCredentialsFile(wd + "/firebaseKey.json")

	//initialize firebase app
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalln("error initializing app: ", err)
	}

	//initialize fcm app
	client, err := app.Messaging(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	//initialize notification service
	fcmRepo := repositories.NewFirebaseMessagingRepository(client)
	fcmService := services.NewNotificationService(fcmRepo)

	//Start Message Broker consumer
	amqpConn := msg_broker.New()
	consumer := msg_broker.NewConsumer(&fcmService, amqpConn)
	consumer.Listen()

}
