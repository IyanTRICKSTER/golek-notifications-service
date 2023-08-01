package main

import (
	"context"
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golek_notifications_service/cmd/msg_broker"
	"golek_notifications_service/pkg/controllers"
	"golek_notifications_service/pkg/database"
	"golek_notifications_service/pkg/models"
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

	//init DB
	db := database.Database{
		Host:     os.Getenv("DB_HOST"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DbName:   os.Getenv("DB_NAME"),
		DbPort:   os.Getenv("DB_PORT"),
	}
	err = db.Connect()
	if err != nil {
		log.Fatalln("error initializing db: ", err)
	}

	db.DropTable(models.Message{})
	db.MigrateTable(models.Message{})

	//load fcm key
	wd, _ := os.Getwd()
	opt := option.WithCredentialsFile(wd + "/firebaseKey.json")

	//initialize firebase app
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalln("error initializing firebase: ", err)
	}

	//initialize fcm app
	client, err := app.Messaging(context.Background())
	if err != nil {
		log.Fatalln("error initializing firebase: ", err)
	}

	//initialize notification service
	fcmRepo := repositories.NewFirebaseMessagingRepository(client)
	notifRepo := repositories.NewNotificationRepo(&db)
	fcmService := services.NewNotificationService(fcmRepo, notifRepo)

	//Start Message Broker consumer
	amqpConn := msg_broker.New()
	consumer := msg_broker.NewConsumer(&fcmService, amqpConn)
	go consumer.Listen()

	//Setup Http Engine
	//Enable Gin Debugging Mode
	//gin.SetMode(gin.ReleaseMode)
	httpEngine := gin.Default()

	controllers.RunNotifController(httpEngine, &fcmService)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "80"
	}

	err = httpEngine.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
