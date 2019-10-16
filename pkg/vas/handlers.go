package vas

import (
	log "git.ronaksoftware.com/blip/server/pkg/logger"
	"git.ronaksoftware.com/blip/server/pkg/user"
	"git.ronaksoftware.com/blip/server/pkg/vas/saba"
	"github.com/kataras/iris"
	"go.uber.org/zap"
)

/*
   Creation Time: 2019 - Oct - 06
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func MCINotification(ctx iris.Context) {
	customerNumber := ctx.URLParam("msisdn")
	status := ctx.URLParam("status")
	amount := ctx.URLParamIntDefault("amount", 0)
	serviceID := ctx.URLParam("serviceId")
	channel := ctx.URLParam("channel")
	dateTime := ctx.URLParamIntDefault("datetime", 0)
	log.Info("VAS NOTIFICATION RECEIVED",
		zap.String("Number", customerNumber),
		zap.String("Status", status),
		zap.Int("Amount", amount),
		zap.String("ServiceID", serviceID),
		zap.String("Channel", channel),
		zap.Int("DateTime", dateTime),
		zap.String("ClientIP", ctx.RemoteAddr()),
	)
	switch status {
	case MciNotificationStatusSubscription:
		u, err := user.GetByPhone(customerNumber)
		if err != nil {
			log.Error("Error On Subscription", zap.Error(err), zap.String("Phone", customerNumber))
			return
		}
		u.VasPaid = true
		err = user.Save(u)
		if err != nil {
			log.Error("Error On Subscription", zap.Error(err), zap.String("Phone", customerNumber))
			return
		}
		res, err := saba.SendMessage(customerNumber, WelcomeMessage)
		if err != nil {
			log.Warn("Error On SendMessage (Subsription)",
				zap.Error(err),
				zap.String("Number", customerNumber),
				zap.String("Status", status),
				zap.Int("Amount", amount),
				zap.String("ServiceID", serviceID),
				zap.String("Channel", channel),
				zap.Int("DateTime", dateTime),
			)
			return
		}
		switch res.StatusCode {
		case saba.SuccessfulCode:
		default:
			log.Info("SendMessage Status",
				zap.String("Status", res.Status),
				zap.String("StatusCode", res.StatusCode),
			)
		}
	case MciNotificationStatusUnsubscription:
		u, err := user.GetByPhone(customerNumber)
		if err != nil {
			log.Error("Error On Subscription", zap.Error(err), zap.String("Phone", customerNumber))
			return
		}
		u.VasPaid = false
		err = user.Save(u)
		if err != nil {
			log.Error("Error On Subscription", zap.Error(err), zap.String("Phone", customerNumber))
			return
		}
		res, err := saba.SendMessage(customerNumber, GoodbyeMessage)
		if err != nil {
			log.Warn("Error On SendMessage (UnSubscription)",
				zap.Error(err),
				zap.String("Number", customerNumber),
				zap.String("Status", status),
				zap.Int("Amount", amount),
				zap.String("ServiceID", serviceID),
				zap.String("Channel", channel),
				zap.Int("DateTime", dateTime),
			)
			return
		}
		switch res.StatusCode {
		case saba.SuccessfulCode:
		default:
			log.Info("SendMessage Status ",
				zap.String("Status", res.Status),
				zap.String("StatusCode", res.StatusCode),
			)
		}
	case MciNotificationStatusActive:
	case MciNotificationStatusDeleted:
	case MciNotificationStatusFailed:
	case MciNotificationStatusPosPaid:
	}
}

func MCIMo(ctx iris.Context) {
	customerNumber := ctx.URLParam("msisdn")
	serviceID := ctx.URLParam("serviceId")
	message := ctx.URLParam("message")

	switch message {
	case "off", "1", "۱", "خاموش":
		u, err := user.GetByPhone(customerNumber)
		if err != nil {
			log.Warn("Unsubscribe received but no user exists", zap.String("Phone", customerNumber))
			return
		}

		res, err := saba.Unsubscribe(customerNumber)
		if err != nil {
			log.Warn("Error On SendMessage (Off Request)",
				zap.Error(err),
				zap.String("Number", customerNumber),
				zap.String("ServiceID", serviceID),
				zap.String("Message", message),
			)
			return
		}
		u.VasPaid = false
		err = user.Save(u)
		if err != nil {
			log.Error("Could not save user", zap.String("UserID", u.ID))
			return
		}
		log.Info("User Unsubscribed",
			zap.String("Status", res),
			zap.String("UserID", u.ID),
			zap.String("Phone", u.Phone),
		)

	case "":
		res, err := saba.SendMessage(customerNumber, EmptyMessage)
		if err != nil {
			log.Warn("Error On SendMessage (EmptyMessage)",
				zap.Error(err),
				zap.String("Number", customerNumber),
				zap.String("ServiceID", serviceID),
				zap.String("Message", message),
			)
			return
		}
		switch res.StatusCode {
		case saba.SuccessfulCode:
		default:
			log.Info("SendMessage Status",
				zap.String("Status", res.Status),
				zap.String("StatusCode", res.StatusCode),
			)
		}
	default:
		res, err := saba.SendMessage(customerNumber, JunkMessage)
		if err != nil {
			log.Warn("Error On SendMessage (JunkMessage)",
				zap.Error(err),
				zap.String("Number", customerNumber),
				zap.String("ServiceID", serviceID),
				zap.String("Message", message),
			)
			return
		}
		switch res.StatusCode {
		case saba.SuccessfulCode:
		default:
			log.Info("SendMessage Status",
				zap.String("Status", res.Status),
				zap.String("StatusCode", res.StatusCode),
			)
		}

	}

}
