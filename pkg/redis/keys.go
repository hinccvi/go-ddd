package redis

type Key string

const (
	Prefix               string = "wsa:"
	RefreshTokenKey      string = "refresh_token:"
	IncorrectPasswordKey string = "incorrect_password:"
	SmsCooldownKey       string = "sms_cooldown:"
	SmsCodeKey           string = "sms_code:"
	SmsLimitKey          string = "sms_limit:"
	SmsAttemptKey        string = "sms_attempt:"
)

func GetRedisKey(key string) string {
	return Prefix + key
}
