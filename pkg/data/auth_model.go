package data

import (
	"time"

	"github.com/abhik-99/passwordless-login/pkg/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Auth struct {
	userId primitive.ObjectID
	otp    string
}

var (
	redisDb = config.Rdb
	rCtx    = config.RedisCtx
)

func (a *Auth) SetOTPForUser() error {
	return redisDb.Set(rCtx, a.userId.Hex(), a.otp, 30*time.Minute).Err()
}

func (a *Auth) CheckOTP() (bool, error) {
	if storedOtp, err := redisDb.Get(rCtx, a.userId.Hex()).Result(); err != nil {
		return false, err
	} else {
		if storedOtp == a.otp {
			return true, nil
		}
	}
	return false, nil
}
