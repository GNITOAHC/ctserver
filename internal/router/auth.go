package router

import (
	"ctserver/internal/authdb"
	"encoding/json"
	"net/http"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type User struct {
	Mail     string `json:"mail"`
	UserName string `json:"username"`
	Phone    string `json:"phone,omitempty"`
	OTP      string `json:"otp"`
}

const (
	OTPPeriod            = 60
	TokenPrefix          = "Bearer "
	AccessTokenDuration  = time.Hour * 12
	RefreshTokenDuration = time.Hour * 24 * 7
)

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
	otp, err := rr.genOTP(u.Mail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
	if u.Mail == "" || u.OTP == "" || u.UserName == "" {
		http.Error(w, "Mail and OTP are required", http.StatusBadRequest)
	}

	if exist, err := rr.helper.CheckUserExist(u.Mail); err != nil || exist {
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, "User already exists", http.StatusBadRequest)
		}
		return
	}
	if exist, err := rr.helper.CheckUsernameExist(u.UserName); err != nil || exist {
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, "Username already used", http.StatusBadRequest)
		}
		return
	}

	// Verify OTP
	valid, err := rr.verifyOTP(u.Mail, u.OTP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !valid {
		http.Error(w, "Invalid OTP", http.StatusUnauthorized)
		return
	}

	// Register user
	err = rr.helper.RegisterUser(u.Mail, u.Phone, u.UserName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("OTP verified"))
	return
}

func (rr *Router) Login(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if u.Mail == "" {
		http.Error(w, "Mail and OTP are required", http.StatusBadRequest)
	}

	// Generate OTP
	otp, err := rr.genOTP(u.Mail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send OTP
	rr.mailer.Send(u.Mail, "OTP Code from ctlink.com", "Your OTP code is: "+otp)

	w.Write([]byte("OTP sent"))
	return
}

type LoginResponse struct {
	RefreshToken string `json:"refreshToken"`
	Username     string `json:"username"`
	Mail         string `json:"mail"`
}

func (rr *Router) LoginVerify(w http.ResponseWriter, r *http.Request) {
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
	valid, err := rr.verifyOTP(u.Mail, u.OTP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !valid {
		http.Error(w, "Invalid OTP", http.StatusUnauthorized)
		return
	}

	// Login user (Generate Access Token & Refresh Token)
	accessToken, err := authdb.NewAccessToken(u.Mail, u.UserName, rr.config.JWTSecret, TokenPrefix, AccessTokenDuration)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	refreshToken, err := rr.authdb.NewSession(authdb.Session{
		Mail:     u.Mail,
		Ip:       r.RemoteAddr,
		Location: r.Header.Get("X-Forwarded-For"),
		Device:   r.Header.Get("User-Agent"),
	}, RefreshTokenDuration, TokenPrefix, rr.config.JWTSecret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the tokens in the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    accessToken,
		MaxAge:   int(AccessTokenDuration),
		Path:     "/",
		HttpOnly: false,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})

	username, err := rr.helper.GetUsername(u.Mail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(&LoginResponse{
		RefreshToken: refreshToken,
		Username:     username,
		Mail:         u.Mail,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}
