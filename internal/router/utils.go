package router

import (
	"errors"
	"time"

	"github.com/pquerna/otp/totp"
)

// genOTP generates a new OTP for the given mail
func (rr *Router) genOTP(mail string) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "ctlink.com",
		AccountName: mail,
		Period:      60, // OTP valid for 60 seconds
	})
	if err != nil {
		return "", err
	}
	otp, err := totp.GenerateCodeCustom(key.Secret(), time.Now(), OtpValidOpts)
	if err != nil {
		return "", err
	}
	rr.cache.Set(mail, key.Secret(), OTPPeriod*time.Second)

	return otp, nil
}

// verifyOTP verifies the OTP for the given mail
func (rr *Router) verifyOTP(mail, otp string) (bool, error) {
	key, ok := rr.cache.Get(mail)
	if !ok {
		return false, errors.New("OTP not found")
	}
	valid, err := totp.ValidateCustom(otp, key, time.Now(), OtpValidOpts)
	if err != nil {
		return false, errors.New("Validation error")
	}
	if !valid {
		return false, nil
	}
	return true, nil
}
