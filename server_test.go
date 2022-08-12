package main

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHandlerLogin(t *testing.T) {
	p, err := bcrypt.GenerateFromPassword([]byte("testwachtwoord"), bcrypt.MinCost)
	if err != nil {
		t.Errorf("%v", err)
	}
	u := user{
		Username: "admin",
		Password: p,
	}
	SaveToJSON(map[string]user{u.Username: u}, fnameUsers)
}
