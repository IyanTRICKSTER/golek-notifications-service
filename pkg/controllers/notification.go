package controllers

import (
	"github.com/gin-gonic/gin"
	"golek_notifications_service/pkg/contracts"
	status "golek_notifications_service/pkg/operation_status"
	"golek_notifications_service/pkg/response"
	"strconv"
)

type NotificationController struct {
	notifSvc contracts.NotificationService
}

func RunNotifController(
	httpEngine *gin.Engine,
	notificationService *contracts.NotificationService,
) {

	notifCtrl := NotificationController{notifSvc: *notificationService}

	v1 := httpEngine.Group("api/v1")
	v1.GET("notifications", notifCtrl.FetchAllNotification)
	v1.GET("notifications/subs/:registrationToken", notifCtrl.NotificationSubscribe)
	v1.GET("notifications/unsubs/:registrationToken", notifCtrl.NotificationUnsubscribe)
}

func (c *NotificationController) FetchAllNotification(ctx *gin.Context) {

	page := ctx.Query("page")
	if page == "" {
		ctx.JSON(400,
			response.New(false, status.Failed, "query page missing", nil).ToMapStringInterface())
		return
	}

	intPage, err := strconv.Atoi(page)
	if err != nil {
		ctx.JSON(400,
			response.New(false, status.Failed, err.Error(), nil).ToMapStringInterface())
		return
	}

	perPage := ctx.Query("per-page")
	if perPage == "" {
		ctx.JSON(400,
			response.New(false, status.Failed, "query per-page missing", nil).ToMapStringInterface())
		return
	}

	intPerPage, err := strconv.Atoi(perPage)
	if err != nil {
		ctx.JSON(400,
			response.New(false, status.Failed, err.Error(), nil).ToMapStringInterface())
		return
	}

	res := c.notifSvc.FetchAll(ctx, uint(intPage), uint(intPerPage))
	if res.IsFailed() {
		ctx.JSON(400, res.ToMapStringInterface())
		return
	}
	ctx.JSON(200, res.ToMapStringInterface())
	return
}

func (c *NotificationController) NotificationSubscribe(ctx *gin.Context) {
	res := c.notifSvc.SubscribeTopic(ctx, ctx.Param("registrationToken"))
	if res.IsFailed() {
		ctx.JSON(400, res.ToMapStringInterface())
		return
	}
	ctx.JSON(200, res.ToMapStringInterface())
	return
}

func (c *NotificationController) NotificationUnsubscribe(ctx *gin.Context) {
	res := c.notifSvc.UnsubscribeTopic(ctx, ctx.Param("registrationToken"))
	if res.IsFailed() {
		ctx.JSON(400, res.ToMapStringInterface())
		return
	}
	ctx.JSON(200, res.ToMapStringInterface())
	return
}
