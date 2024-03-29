package repositories

import (
	"context"
	"encoding/base64"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"golek_notifications_service/pkg/contracts"
	"golek_notifications_service/pkg/models"
	status "golek_notifications_service/pkg/operation_status"
	"log"
	"os"
	"time"
)

type FirebaseMessagingRepository struct {
	client *messaging.Client
}

func NewFirebaseMessagingRepository(client *messaging.Client) contracts.ICloudMessageRepo {
	return &FirebaseMessagingRepository{client: client}
}

func (f *FirebaseMessagingRepository) SendToTopic(ctx context.Context, topic string, message *models.Message) (status.OperationStatus, string, error) {

	fcmMessage := messaging.Message{
		Data: map[string]string{
			"title":     message.Title,
			"body":      message.Body,
			"image_url": message.ImageUrl,
		},
		Notification: &messaging.Notification{
			Title:    message.Title,
			Body:     message.Body,
			ImageURL: message.ImageUrl,
		},
		Topic: topic,
	}

	response, err := f.client.Send(ctx, &fcmMessage)
	if err != nil {
		log.Printf("ERROR FCM Repository:SendToTopic > %v", err.Error())
		return status.SendMessageFailed, "", err
	}

	return status.SendMessageSuccess, response, nil
}

func (f *FirebaseMessagingRepository) SubscribeToTopic(ctx context.Context, topic string, registrationTokens []string) (status.OperationStatus, string, error) {
	//topic := "highScores"

	// [START subscribe_golang]
	// These registration tokens come from the client FCM SDKs.
	//registrationTokens := []string{
	//	"YOUR_REGISTRATION_TOKEN_1",
	//	// ...
	//	"YOUR_REGISTRATION_TOKEN_n",
	//}

	// Subscribe the devices corresponding to the registration tokens to the
	// topic.
	response, err := f.client.SubscribeToTopic(ctx, registrationTokens, topic)
	if err != nil {
		log.Println(err)
		return status.Failed, "", err
	}
	// See the TopicManagementResponse reference documentation
	// for the contents of response.
	fmt.Println(response.SuccessCount, "tokens were subscribed successfully")
	// [END subscribe_golang]
	return status.Success, "", nil
}

func (f *FirebaseMessagingRepository) UnsubscribeFromTopic(ctx context.Context, topic string, registrationTokens []string) (status.OperationStatus, string, error) {
	//topic := "highScores"

	// [START unsubscribe_golang]
	// These registration tokens come from the client FCM SDKs.
	//registrationTokens := []string{
	//	"YOUR_REGISTRATION_TOKEN_1",
	//	// ...
	//	"YOUR_REGISTRATION_TOKEN_n",
	//}

	// Unsubscribe the devices corresponding to the registration tokens from
	// the topic.
	response, err := f.client.UnsubscribeFromTopic(ctx, registrationTokens, topic)
	if err != nil {
		log.Println(err)
		return status.Failed, "", err
	}
	// See the TopicManagementResponse reference documentation
	// for the contents of response.
	fmt.Println(response.SuccessCount, "tokens were unsubscribed successfully")
	// [END unsubscribe_golang]
	return status.Success, "", nil
}

func getDecodedFireBaseKey() ([]byte, error) {

	fireBaseAuthKey := os.Getenv("FIREBASE_AUTH_KEY")

	decodedKey, err := base64.StdEncoding.DecodeString(fireBaseAuthKey)
	if err != nil {
		return nil, err
	}

	return decodedKey, nil
}

func sendToToken(app *firebase.App) {

	// [START send_to_token_golang]
	// Obtain a messaging.Client from the App.
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	// This registration token comes from the client FCM SDKs.
	registrationToken := "YOUR_REGISTRATION_TOKEN"

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
		Token: registrationToken,
	}

	// Send a message to the device corresponding to the provided
	// registration token.
	response, err := client.Send(ctx, message)
	if err != nil {
		log.Fatalln(err)
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", response)
	// [END send_to_token_golang]
}

func sendToTopic(ctx context.Context, client *messaging.Client) {
	// [START send_to_topic_golang]
	// The topic name can be optionally prefixed with "/topics/".
	topic := "newPost"

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
		Topic: topic,
		Notification: &messaging.Notification{
			Title:    "Message from backend",
			Body:     "Mantap!!!",
			ImageURL: "https://i.pinimg.com/originals/dd/c3/af/ddc3afe873e140699868ecef9419caa3.png",
		},
	}

	// Send a message to the devices subscribed to the provided topic.
	response, err := client.Send(ctx, message)
	if err != nil {
		log.Fatalln(err)
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", response)
	// [END send_to_topic_golang]
}

func sendToCondition(ctx context.Context, client *messaging.Client) {
	// [START send_to_condition_golang]
	// Define a condition which will send to devices which are subscribed
	// to either the Google stock or the tech industry topics.
	condition := "'stock-GOOG' in topics || 'industry-tech' in topics"

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
		Condition: condition,
	}

	// Send a message to devices subscribed to the combination of topics
	// specified by the provided condition.
	response, err := client.Send(ctx, message)
	if err != nil {
		log.Fatalln(err)
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", response)
	// [END send_to_condition_golang]
}

func sendAll(ctx context.Context, client *messaging.Client) {
	// This registration token comes from the client FCM SDKs.
	registrationToken := "YOUR_REGISTRATION_TOKEN"

	// [START send_all]
	// Create a list containing up to 500 messages.
	messages := []*messaging.Message{
		{
			Notification: &messaging.Notification{
				Title: "Price drop",
				Body:  "5% off all electronics",
			},
			Token: registrationToken,
		},
		{
			Notification: &messaging.Notification{
				Title: "Price drop",
				Body:  "2% off all books",
			},
			Topic: "readers-club",
		},
	}

	br, err := client.SendAll(context.Background(), messages)
	if err != nil {
		log.Fatalln(err)
	}

	// See the BatchResponse reference documentation
	// for the contents of response.
	fmt.Printf("%d messages were sent successfully\n", br.SuccessCount)
	// [END send_all]
}

func sendMulticast(ctx context.Context, client *messaging.Client) {
	// [START send_multicast]
	// Create a list containing up to 500 registration tokens.
	// This registration tokens come from the client FCM SDKs.
	registrationTokens := []string{
		"YOUR_REGISTRATION_TOKEN_1",
		// ...
		"YOUR_REGISTRATION_TOKEN_n",
	}
	message := &messaging.MulticastMessage{
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
		Tokens: registrationTokens,
	}

	br, err := client.SendMulticast(context.Background(), message)
	if err != nil {
		log.Fatalln(err)
	}

	// See the BatchResponse reference documentation
	// for the contents of response.
	fmt.Printf("%d messages were sent successfully\n", br.SuccessCount)
	// [END send_multicast]
}

func sendMulticastAndHandleErrors(ctx context.Context, client *messaging.Client) {
	// [START send_multicast_error]
	// Create a list containing up to 500 registration tokens.
	// This registration tokens come from the client FCM SDKs.
	registrationTokens := []string{
		"YOUR_REGISTRATION_TOKEN_1",
		// ...
		"YOUR_REGISTRATION_TOKEN_n",
	}
	message := &messaging.MulticastMessage{
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
		Tokens: registrationTokens,
	}

	br, err := client.SendMulticast(context.Background(), message)
	if err != nil {
		log.Fatalln(err)
	}

	if br.FailureCount > 0 {
		var failedTokens []string
		for idx, resp := range br.Responses {
			if !resp.Success {
				// The order of responses corresponds to the order of the registration tokens.
				failedTokens = append(failedTokens, registrationTokens[idx])
			}
		}

		fmt.Printf("List of tokens that caused failures: %v\n", failedTokens)
	}
	// [END send_multicast_error]
}

func sendDryRun(ctx context.Context, client *messaging.Client) {
	message := &messaging.Message{
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
		Token: "token",
	}

	// [START send_dry_run_golang]
	// Send a message in the dry run mode.
	response, err := client.SendDryRun(ctx, message)
	if err != nil {
		log.Fatalln(err)
	}
	// Response is a message ID string.
	fmt.Println("Dry run successful:", response)
	// [END send_dry_run_golang]
}

func androidMessage() *messaging.Message {
	// [START android_message_golang]
	oneHour := time.Duration(1) * time.Hour
	message := &messaging.Message{
		Android: &messaging.AndroidConfig{
			TTL:      &oneHour,
			Priority: "normal",
			Notification: &messaging.AndroidNotification{
				Title: "$GOOG up 1.43% on the day",
				Body:  "$GOOG gained 11.80 points to close at 835.67, up 1.43% on the day.",
				Icon:  "stock_ticker_update",
				Color: "#f45342",
			},
		},
		Topic: "industry-tech",
	}
	// [END android_message_golang]
	return message
}

func apnsMessage() *messaging.Message {
	// [START apns_message_golang]
	badge := 42
	message := &messaging.Message{
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-priority": "10",
			},
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Alert: &messaging.ApsAlert{
						Title: "$GOOG up 1.43% on the day",
						Body:  "$GOOG gained 11.80 points to close at 835.67, up 1.43% on the day.",
					},
					Badge: &badge,
				},
			},
		},
		Topic: "industry-tech",
	}
	// [END apns_message_golang]
	return message
}

func webpushMessage() *messaging.Message {
	// [START webpush_message_golang]
	message := &messaging.Message{
		Webpush: &messaging.WebpushConfig{
			Notification: &messaging.WebpushNotification{
				Title: "$GOOG up 1.43% on the day",
				Body:  "$GOOG gained 11.80 points to close at 835.67, up 1.43% on the day.",
				Icon:  "https://my-server/icon.png",
			},
		},
		Topic: "industry-tech",
	}
	// [END webpush_message_golang]
	return message
}

func allPlatformsMessage() *messaging.Message {
	// [START multi_platforms_message_golang]
	oneHour := time.Duration(1) * time.Hour
	badge := 42
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: "$GOOG up 1.43% on the day",
			Body:  "$GOOG gained 11.80 points to close at 835.67, up 1.43% on the day.",
		},
		Android: &messaging.AndroidConfig{
			TTL: &oneHour,
			Notification: &messaging.AndroidNotification{
				Icon:  "stock_ticker_update",
				Color: "#f45342",
			},
		},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Badge: &badge,
				},
			},
		},
		Topic: "industry-tech",
	}
	// [END multi_platforms_message_golang]
	return message
}
