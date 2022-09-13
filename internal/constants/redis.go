package constants

type RedisKey string

const (
	prefix               RedisKey = "app:"
	RefreshTokenKey      RedisKey = "refresh_token:"
	IncorrectPasswordKey RedisKey = "incorrect_password:"
	SmsCooldownKey       RedisKey = "sms_cooldown:"
	SmsCodeKey           RedisKey = "sms_code:"
	SmsLimitKey          RedisKey = "sms_limit:"
	SmsAttemptKey        RedisKey = "sms_attempt:"
)

func GetRedisKey(key RedisKey) string {
	return string(prefix) + string(key)
}
