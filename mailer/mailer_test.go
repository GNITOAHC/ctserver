package mailer

import (
	"ctserver/internal/config"
	"testing"
)

func TestSend(t *testing.T) {
	c, _ := config.New("../.env")
	m := New(c.SMTPFrom, c.SMTPPass, c.SMTPHost, c.SMTPPort)
	err := m.Send(c.TestMail, "Test subject", "<p style='color: blue'>Test body</p>")
	if err != nil {
		t.Fatal(err)
	}
}
