package contracts

import (
	"context"
	"golek_notifications_service/pkg/models"
	status "golek_notifications_service/pkg/operation_status"
)

type INotificationRepository interface {
	FetchAll(ctx context.Context, page uint, perPage uint) ([]models.Message, uint, status.OperationStatus, error)
	Create(ctx context.Context, payload models.Message) (status.OperationStatus, error)
}
