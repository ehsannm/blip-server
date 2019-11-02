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
به سرويس موزيكچي خوش آمديد
لينكدانلود اپليكيشن:
https://getmusicchi.ir/musicchi.apk
براي لغو سرويس كليدواژه off يا خاموش را به همين سرشماره ارسال كنيد
هزينه روزانه سرويس ٦٠٠ تومان
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
