package contracts

import (
	"context"
	"golek_notifications_service/pkg/models"
	"golek_notifications_service/pkg/operation_status"
)

type NotificationRepository interface {
	SendToTopic(ctx context.Context, topic string, message *models.Message) (opStatus status.OperationStatus, response string, err error)
}

type NotificationService interface {
	Broadcast(ctx context.Context, topic string, message *models.Message) (opStatus status.OperationStatus, err error)
}
