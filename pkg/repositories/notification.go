package repositories

import (
	"context"
	"errors"
	"golek_notifications_service/pkg/contracts"
	"golek_notifications_service/pkg/database"
	"golek_notifications_service/pkg/models"
	status "golek_notifications_service/pkg/operation_status"
	"gorm.io/gorm"
)

type NotificationRepo struct {
	db *database.Database
}

func (n *NotificationRepo) FetchAll(ctx context.Context, page uint, perPage uint) ([]models.Message, uint, status.OperationStatus, error) {
	var notifs []models.Message
	var notifTotal int64

	offset := (page - 1) * perPage
	err := n.db.GetConnection().WithContext(ctx).Limit(int(perPage)).Offset(int(offset)).Find(&notifs).Error
	err = n.db.GetConnection().WithContext(ctx).Model(models.Message{}).Count(&notifTotal).Error
	if err != nil {
		return notifs, 0, status.Failed, err
	}
	return notifs, uint(notifTotal), status.Success, nil
}

func (n *NotificationRepo) Create(ctx context.Context, payload models.Message) (status.OperationStatus, error) {
	err := n.db.GetConnection().WithContext(ctx).Create(payload).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return status.ErrorDuplicatedModel, err
		}
		return status.Failed, err
	}
	return status.Success, nil
}

func NewNotificationRepo(db *database.Database) contracts.INotificationRepository {
	return &NotificationRepo{db: db}
}
