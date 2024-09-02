package authdb

import (
	"ctserver/internal/config"
	"testing"
	"time"
)

const token = "Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJNYWlsIjoidGVzdEB0ZXN0LmNvbSIsIlNlc3Npb25Ub2tlbiI6IjQ2OTUwZGI4LWRiOGUtNDU4NC05MzcxLWIyZjAxNTFlMzMxNiIsImV4cCI6MTcyNTI5MjI4OCwiaWF0IjoxNzI1MjkwNDg4fQ.VOclqINQ2T2vPdF8fYWccx8-FZSByDWqzfwhbyJw843ZZbpVEGlMPGFnlbl9ixj68fjvFLZCcjge1_KWGj6Trw"

func TestNewSession(t *testing.T) {
	c, _ := config.New("../../.env")
	db := New(c.AuthDBURI, c.AuthDBName, c.AuthDBCollection, c.JWTSecret)
	session := Session{
		Mail: "test@test.com",
	}
	token, err := db.NewSession(session, 30*time.Minute, "Bearer ", c.JWTSecret)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token)
}

func TestListSession(t *testing.T) {
	c, _ := config.New("../../.env")
	db := New(c.AuthDBURI, c.AuthDBName, c.AuthDBCollection, c.JWTSecret)
	sessions, err := db.ListSessions("test@test.com")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sessions)
}

func TestDeleteSession(t *testing.T) {
	c, _ := config.New("../../.env")
	db := New(c.AuthDBURI, c.AuthDBName, c.AuthDBCollection, c.JWTSecret)
	// session := Session{
	// 	Mail: "test@test.com",
	// }
	// token, err := db.NewSession(session, 30*time.Second, "Bearer ", c.JWTSecret)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// token := "Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJNYWlsIjoidGVzdEB0ZXN0LmNvbSIsIlNlc3Npb25Ub2tlbiI6Ijg5ODY2ZTM2LTZmOGUtNDZhMS04Nzk1LTdiNDg5MWMwNTM2OCIsImV4cCI6MTcyNTI4OTQ5NCwiaWF0IjoxNzI1Mjg5NDY0fQ.ZSkf3uHTDMN_EG7JjmGTRMjqj6pOGxeKBgpDEDc5RWkb8m06nzCsWSP3-Hpm4Le0J-99w_62qH6x1ezvnw0Yaw"
	err := db.DeleteSession(token)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheckSession(t *testing.T) {
	c, _ := config.New("../../.env")
	db := New(c.AuthDBURI, c.AuthDBName, c.AuthDBCollection, c.JWTSecret)
	// session := Session{
	// 	Mail: "test@test.com",
	// }
	// token, err := db.NewSession(session, 30*time.Second, "Bearer ", c.JWTSecret)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// token := "Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJNYWlsIjoidGVzdEB0ZXN0LmNvbSIsIlNlc3Npb25Ub2tlbiI6Ijg5ODY2ZTM2LTZmOGUtNDZhMS04Nzk1LTdiNDg5MWMwNTM2OCIsImV4cCI6MTcyNTI4OTQ5NCwiaWF0IjoxNzI1Mjg5NDY0fQ.ZSkf3uHTDMN_EG7JjmGTRMjqj6pOGxeKBgpDEDc5RWkb8m06nzCsWSP3-Hpm4Le0J-99w_62qH6x1ezvnw0Yaw"
	// decoded, err := Decode(token, c.JWTSecret, "Bearer ")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(decoded)

	// decoded, err := jwt.Parse(token, c.JWTSecret, "Bearer ")
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// if isExpired, err := jwt.IsExpired(decoded); isExpired || err != nil {
	// 	t.Fatal(err)
	// }
	// log.Print(decoded)

	valid, err := db.ValidSession(token, c.JWTSecret, "Bearer ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(valid)
}
