package services

import (
	"context"
	"golek_notifications_service/pkg/contracts"
	"golek_notifications_service/pkg/models"
	status "golek_notifications_service/pkg/operation_status"
	"golek_notifications_service/pkg/response"
	"log"
)

type NotificationService struct {
	cloudMessageRepo contracts.ICloudMessageRepo
	notifRepo        contracts.INotificationRepository
}

func (n NotificationService) SubscribeTopic(ctx context.Context, registrationToken string) response.Response {
	opStatus, _, err := n.cloudMessageRepo.SubscribeToTopic(ctx, "newPost", []string{registrationToken})
	if err != nil {
		return response.New(false, opStatus, err.Error(), nil)
	}
	return response.New(true, opStatus, "ok", nil)
}

func (n NotificationService) UnsubscribeTopic(ctx context.Context, registrationToken string) response.Response {
	opStatus, _, err := n.cloudMessageRepo.UnsubscribeFromTopic(ctx, "newPost", []string{registrationToken})
	if err != nil {
		return response.New(false, opStatus, err.Error(), nil)
	}
	return response.New(true, opStatus, "ok", nil)
}

func (n NotificationService) FetchAll(ctx context.Context, page uint, perPage uint) response.Response {
	notifs, totalNotif, opStatus, err := n.notifRepo.FetchAll(ctx, page, perPage)
	if err != nil {
		return response.New(false, opStatus, err.Error(), nil)
	}
	return response.New(true, opStatus, "ok",
		map[string]interface{}{"notifications": notifs, "page": page, "perPage": perPage, "totalData": totalNotif})
}

func (n NotificationService) Broadcast(ctx context.Context, topic string, message *models.Message) (opStatus status.OperationStatus, err error) {

	operationStatus, res, err := n.cloudMessageRepo.SendToTopic(ctx, topic, message)
	if err != nil {
		log.Printf("ERROR Notification Service:Broadcast > %v", err.Error())
		//return operationStatus, err
	}

	opStatus, err = n.notifRepo.Create(ctx, *message)
	if err != nil {
		log.Printf("ERROR Notification Service:Broadcast > %v", err.Error())
	}

	log.Printf("Notification Service:Broadcast > %v", res)
	return operationStatus, nil
}

func NewNotificationService(
	notificationRepo contracts.ICloudMessageRepo,
	notifRepo contracts.INotificationRepository,
) contracts.NotificationService {
	return &NotificationService{
		cloudMessageRepo: notificationRepo,
		notifRepo:        notifRepo,
	}
}
