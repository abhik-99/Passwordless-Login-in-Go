package data

import (
	"time"

	"github.com/abhik-99/passwordless-login/pkg/config"
)

type Auth struct {
	UserId string
	Otp    string
}

var (
	redisDb = config.Rdb
	rCtx    = config.RedisCtx
)

func (a *Auth) SetOTPForUser() error {
	return redisDb.Set(rCtx, a.UserId, a.Otp, 30*time.Minute).Err()
}

func (a *Auth) CheckOTP() (bool, error) {
	if storedOtp, err := redisDb.Get(rCtx, a.UserId).Result(); err != nil {
		return false, err
	} else {
		if storedOtp == a.Otp {
			return true, nil
		}
	}
	return false, nil
}
