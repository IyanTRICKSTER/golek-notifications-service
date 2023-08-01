package contracts

import (
	"context"
	"golek_notifications_service/pkg/models"
	"golek_notifications_service/pkg/operation_status"
	"golek_notifications_service/pkg/response"
)

type ICloudMessageRepo interface {
	SendToTopic(ctx context.Context, topic string, message *models.Message) (opStatus status.OperationStatus, response string, err error)
	SubscribeToTopic(ctx context.Context, topic string, registrationTokens []string) (opStatus status.OperationStatus, response string, err error)
	UnsubscribeFromTopic(ctx context.Context, topic string, registrationTokens []string) (opStatus status.OperationStatus, response string, err error)
}

type NotificationService interface {
	Broadcast(ctx context.Context, topic string, message *models.Message) (opStatus status.OperationStatus, err error)
	FetchAll(ctx context.Context, page uint, perPage uint) response.Response
	SubscribeTopic(ctx context.Context, registrationToken string) response.Response
	UnsubscribeTopic(ctx context.Context, registrationToken string) response.Response
}
