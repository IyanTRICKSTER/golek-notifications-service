package repositories

import (
	"context"
	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"golek_notifications_service/pkg/models"
	"google.golang.org/api/option"
	"os"
	"testing"
)

func TestFirebaseCloudMessagingRepository(t *testing.T) {

	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal(err)
	}

	// Initialize firebase app
	wd, _ := os.Getwd()
	opt := option.WithCredentialsFile(wd + "/../../firebaseKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		t.Fatal("error initializing app: ", err)
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	notifyRepo := NewFirebaseMessagingRepository(client)
	_, res, err := notifyRepo.SendToTopic(context.Background(), "newPost", &models.Message{
		Title:    "Heya boy!",
		Body:     "Got something?",
		ImageUrl: "https://image.apktoy.com/img/13/com.nkl.xnxx.app/icon.png",
		UserID:   100,
	})
	if err != nil {
		t.Error(res)
		t.Error(err)
	}

	t.Log(res)

}
