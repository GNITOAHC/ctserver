package router

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type User struct {
	Mail  string `json:"mail"`
	Phone string `json:"phone,omitempty"`
	OTP   string `json:"otp"`
}

const OTPPeriod = 60

var OtpValidOpts = totp.ValidateOpts{Period: OTPPeriod, Digits: otp.DigitsSix}

func (rr *Router) Register(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if u.Mail == "" {
		http.Error(w, "Mail is required", http.StatusBadRequest)
	}
	if exist, err := rr.helper.CheckUserExist(u.Mail); err != nil || exist {
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, "User already exists", http.StatusBadRequest)
		}
		return
	}

	// Generate OTP
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "ctlink.com",
		AccountName: u.Mail,
		Period:      60, // OTP valid for 60 seconds
	})
	otp, err := totp.GenerateCodeCustom(key.Secret(), time.Now(), OtpValidOpts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	rr.cache.Set(u.Mail, key.Secret(), OTPPeriod*time.Second)

	// Send OTP
	rr.mailer.Send(u.Mail, "OTP Code from ctlink.com", "Your OTP code is: "+otp)

	w.Write([]byte("OTP sent"))
	return
}

func (rr *Router) RegVerify(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if u.Mail == "" || u.OTP == "" {
		http.Error(w, "Mail and OTP are required", http.StatusBadRequest)
	}

	// Verify OTP
	key, ok := rr.cache.Get(u.Mail)
	if !ok {
		http.Error(w, "OTP not found", http.StatusBadRequest)
		return
	}
	valid, err := totp.ValidateCustom(u.OTP, key, time.Now(), OtpValidOpts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if !valid {
		http.Error(w, "Invalid OTP", http.StatusUnauthorized)
		return
	}

	// Register user
	err = rr.helper.RegisterUser(u.Mail, u.Phone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("OTP verified"))
	return
}
