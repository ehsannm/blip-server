package msg

/*
   Creation Time: 2019 - Sep - 21
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

type Item string

var (
	ErrAccessTokenInvalid      Item = "ACCESS_KEY_INVALID"
	ErrAccessTokenExpired      Item = "ACCESS_KEY_EXPIRED"
	ErrNoPermission            Item = "NO_PERMISSION"
	ErrWriteToDb               Item = "WRITE_TO_DB"
	ErrTokenInvalid            Item = "TOKEN_INVALID"
	ErrTokenExpired            Item = "TOKEN_EXPIRED"
	ErrPermissionIsNotSet      Item = "PERMISSION_NOT_SET"
	ErrPhoneNotValid           Item = "PHONE_NOT_VALID"
	ErrPhoneCodeNotValid       Item = "PHONE_CODE_NOT_VALID"
	ErrPhoneCodeHashNotValid   Item = "PHONE_CODE_HASH_NOT_VALID"
	ErrPeriodNotValid          Item = "PERIOD_NOT_VALID"
	ErrWriteToCache            Item = "WRITE_TO_CACHE"
	ErrReadFromCache           Item = "READ_FROM_CACHE"
	ErrReadFromDb              Item = "READ_FROM_DB"
	ErrCannotUnmarshalRequest  Item = "CANNOT_UNMARSHAL_JSON"
	ErrAlreadyRegistered       Item = "ALREADY_REGISTERED"
	ErrUsernameFormat          Item = "USERNAME_FORMAT"
	ErrUnsupportedCarrier      Item = "UNSUPPORTED_CARRIER"
	ErrSessionInvalid          Item = "SESSION_INVALID"
	ErrBadSoundFile            Item = "BAD_SOUND_FILE"
	ErrNoResponseFromVAS       Item = "NO_RESPONSE_FROM_VAS"
	ErrNoResponseFromSmsServer Item = "NO_RESPONSE_FROM_SMS_PANEL"
	ErrVasIsNotEnabled         Item = "VAS_IS_DISABLED"
	ErrUserNotFound            Item = "USER_NOT_FOUND"
	Err3rdParty                Item = "3RD_PARTY_INVALID_RESPONSE"
	ErrCorruptData             Item = "CORRUPT_DATA"
	ErrLocalIndexFailure       Item = "LOCAL_INDEX_FAILED"
	ErrAlreadyServed           Item = "ALREADY_SERVED"
)
