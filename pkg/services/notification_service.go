package services

import (
	"context"
	"golek_notifications_service/pkg/contracts"
	"golek_notifications_service/pkg/models"
	status "golek_notifications_service/pkg/operation_status"
	"log"
)

type NotificationService struct {
	notificationRepo contracts.NotificationRepository
}

func (n NotificationService) Broadcast(ctx context.Context, topic string, message *models.Message) (opStatus status.OperationStatus, err error) {

	operationStatus, res, err := n.notificationRepo.SendToTopic(ctx, topic, message)
	if err != nil {
		log.Printf("ERROR Notification Service:Broadcast > %v", err.Error())
		return operationStatus, err
	}

	log.Printf("Notification Service:Broadcast > %v", res)
	return operationStatus, nil
}

func NewNotificationService(notificationRepo contracts.NotificationRepository) contracts.NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
	}
}
