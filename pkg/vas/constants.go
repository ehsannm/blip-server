package vas

/*
   Creation Time: 2019 - Oct - 06
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

const (
	MciNotificationStatusSubscription   = "subscription"
	MciNotificationStatusUnsubscription = "unSubscription"
	MciNotificationStatusActive         = "active"
	MciNotificationStatusFailed         = "failed"
	MciNotificationStatusDeleted        = "deleted"
	MciNotificationStatusPosPaid        = "postPaid"
)

const (
	WelcomeMessage = `
مشترک گرامی 
سرویس موزیکچی
با تعرفه روزانه 600 تومان برای شما فعال گردید
جهت غیر فعال سازی میتوانید off یا خاموش را به همین شماره ارسال فرمایید
`
	GoodbyeMessage = `
سرویس موزیکچی با موفقیت غیرفعالگردید
برای فعال سازی مجدد  سرویس موزیکچی عدد 1 را به همین سرشماره ارسال فرمایید
`
	EmptyMessage = `
مشترک گرامی کلید واژه ارسالی نادرست است
برای فعال سازی سرویس موزیکچی با تعرفه روزانه 600 تومان، عدد 1 را به همین سرشمار ارسال فرمایید
`
	JunkMessage = `
مشترک گرامی کلید واژه ارسالی نادرست است
برای فعال سازی سرویس موزیکچی با تعرفه روزانه 600 تومان، عدد 1 را به همین سرشمار ارسال فرمایید
`
)
