package vas

import (
	"fmt"
	log "git.ronaksoftware.com/blip/server/pkg/logger"
	"git.ronaksoftware.com/blip/server/pkg/session"
	"git.ronaksoftware.com/blip/server/pkg/user"
	"git.ronaksoftware.com/blip/server/pkg/vas/saba"
	ronak "git.ronaksoftware.com/ronak/toolbox"
	"github.com/kataras/iris"
	"go.uber.org/zap"
	"strings"
	"time"
)

/*
   Creation Time: 2019 - Oct - 06
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/



type mciNotificationParams struct {
	CustomerNumber string `bson:"number"`
	Status         string `bson:"status"`
	Amount         int    `bson:"amount"`
	ServiceID      string `bson:"service_id"`
	Channel        string `bson:"channel"`
	DateTime       int    `bson:"created_on"`
}

func MCINotification(ctx iris.Context) {
	params := &mciNotificationParams{
		CustomerNumber: ctx.URLParam("msisdn"),
		Status:         ctx.URLParam("status"),
		Amount:         ctx.URLParamIntDefault("amount", 0),
		ServiceID:      ctx.URLParam("serviceId"),
		Channel:        ctx.URLParam("channel"),
		DateTime:       ctx.URLParamIntDefault("datetime", 0),
	}

	log.Info("VAS NOTIFICATION RECEIVED",
		zap.String("Number", params.CustomerNumber),
		zap.String("Status", params.Status),
		zap.String("Channel", params.Channel),
	)
	writeToDB.Enter(nil, params)
	switch params.Status {
	case MciNotificationStatusSubscription:
		subscribe(params)
	case MciNotificationStatusUnsubscription:
		unsubscribe(params)
	case MciNotificationStatusActive:
		// Some who charged
	case MciNotificationStatusDeleted:
		//
	case MciNotificationStatusFailed:
		// Failed could not charge
	case MciNotificationStatusPostPaid:

	}
}
func subscribe(params *mciNotificationParams) {
	u, err := user.GetByPhone(params.CustomerNumber)
	if err != nil {
		userID := fmt.Sprintf("U%s", ronak.RandomID(32))
		timeNow := time.Now().Unix()
		u = &user.User{
			ID:        userID,
			Username:  fmt.Sprintf("USER%s", strings.ToUpper(ronak.RandomID(12))),
			Phone:     params.CustomerNumber,
			Email:     "",
			CreatedOn: timeNow,
			Disabled:  false,
		}
	}
	u.VasPaid = true
	err = user.Save(u)
	if err != nil {
		log.Error("Error On Subscription", zap.Error(err), zap.String("Phone", params.CustomerNumber))
		return
	}
	res, err := saba.SendMessage(params.CustomerNumber, WelcomeMessage)
	if err != nil {
		log.Warn("Error On SendMessage (Subsription)",
			zap.Error(err),
			zap.String("Number", params.CustomerNumber),
			zap.String("Status", params.Status),
			zap.Int("Amount", params.Amount),
			zap.String("ServiceID", params.ServiceID),
			zap.String("Channel", params.Channel),
			zap.Int("DateTime", params.DateTime),
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
func unsubscribe(params *mciNotificationParams) {
	u, err := user.GetByPhone(params.CustomerNumber)
	if err != nil {
		log.Error("Error On Unsubscription", zap.Error(err), zap.String("Phone", params.CustomerNumber))
		return
	}
	u.VasPaid = false
	err = user.Save(u)
	if err != nil {
		log.Error("Error On Subscription (Update User)", zap.Error(err), zap.String("Phone", params.CustomerNumber))
		return
	}
	err = session.RemoveAll(u.ID)
	if err != nil {
		log.Error("Error On Subscription (Remove Sessions)", zap.Error(err), zap.String("Phone", params.CustomerNumber))
		return
	}
	res, err := saba.SendMessage(params.CustomerNumber, GoodbyeMessage)
	if err != nil {
		log.Warn("Error On SendMessage (Unsubscription)",
			zap.Error(err),
			zap.String("Number", params.CustomerNumber),
			zap.String("Status", params.Status),
			zap.Int("Amount", params.Amount),
			zap.String("ServiceID", params.ServiceID),
			zap.String("Channel", params.Channel),
			zap.Int("DateTime", params.DateTime),
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
