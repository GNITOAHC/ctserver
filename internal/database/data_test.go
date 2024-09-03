package database

import (
	"ctserver/internal/config"
	"testing"
)

func TestNewUrlDefault(t *testing.T) {
	c, err := config.New("../../.env")
	if err != nil {
		t.Fatal(err)
	}
	db := New(c.DatabaseURI)
	shortened, err := db.NewUrlDefault(Data{Content: "https://www.google.com", Path: "_universal2", Duration: -1})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(shortened)
}

func TestGetUrlDefault(t *testing.T) {
	c, err := config.New("../../.env")
	if err != nil {
		t.Fatal(err)
	}
	db := New(c.DatabaseURI)
	original, err := db.GetUrlDefault("_universal1")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(original)
}

func TestNewUrlCustom(t *testing.T) {
	c, err := config.New("../../.env")
	if err != nil {
		t.Fatal(err)
	}
	db := New(c.DatabaseURI)
	shortened, err := db.NewUrlCustom(Data{
		Username: "wowctchen",
		Content:  "https://www.google.com",
		Path:     "/custompath2",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(shortened)
}

func TestGetUrlCustom(t *testing.T) {
	c, err := config.New("../../.env")
	if err != nil {
		t.Fatal(err)
	}
	db := New(c.DatabaseURI)
	original, err := db.GetUrlCustom("wowctchen", "/custompath1")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(original)
}
