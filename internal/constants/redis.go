package constants

import "time"

type RedisKey string

const (
	MaxLoginAttempt = 5

	IncorrectPasswordExpiration = 24 * time.Hour

	Prefix               RedisKey = "app:"
	RefreshTokenKey      RedisKey = "refresh_token:"
	IncorrectPasswordKey RedisKey = "incorrect_password:"
	SmsCooldownKey       RedisKey = "sms_cooldown:"
	SmsCodeKey           RedisKey = "sms_code:"
	SmsLimitKey          RedisKey = "sms_limit:"
	SmsAttemptKey        RedisKey = "sms_attempt:"
)

func GetRedisKey(key RedisKey) RedisKey {
	return Prefix + key
}
