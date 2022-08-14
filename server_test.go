package main

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHandlerLogin(t *testing.T) {
	p, err := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.MinCost)
	if err != nil {
		t.Errorf("%v", err)
	}
	u := user{
		Username: "admin",
		Password: p,
	}
	SaveToJSON(map[string]user{u.Username: u}, fnameUsers)
}

func TestStartsWith(t *testing.T) {
	cases := []struct {
		s      string
		substr string
		want   bool
	}{
		{"http://ditiseentest.nl", "http", true},
		{"https://ditiseentest.nl", "http", true},
		{"www.ditiseentest.nl", "http", false},
		{"www.ditiseentest.nl", "www", true},
		{"HTTP://ditiseentest.nl", "http", true},
		{"HTTPS://ditiseentest.nl", "https", true},
		{"WWW.ditiseentest.nl", "www", true},
		{"https://www.leukerecepten.nl/snelle-plantaardige-kip-tandoori-recept-kom-gratis-proeven/", "https:", true},
	}
	for i, c := range cases {
		if got := startsWith(c.s, c.substr); got != c.want {
			t.Errorf("Case %v failed: '%v' '%v'. Want: %v Got: %v", i, c.s, c.substr, c.want, got)
		}
	}
}

func TestRemoveFromSlice(t *testing.T) {
	cases := []struct {
		xi   []string
		i    int
		want []string
	}{
		{[]string{"a", "b", "c", "d", "e", "f", "g", "h"}, 3, []string{"a", "b", "c", "e", "f", "g", "h"}},
		{[]string{"a", "b", "c", "d", "e", "f", "g", "h"}, 0, []string{"b", "c", "d", "e", "f", "g", "h"}},
		{[]string{"a", "b", "c", "d", "e", "f", "g", "h"}, 7, []string{"a", "b", "c", "d", "e", "f", "g"}},
	}
	for i, c := range cases {
		got := removeFromSlice(c.xi, c.i)
		for j := 0; j < len(got); j++ {
			if c.want[j] != got[j] {
				t.Errorf("Case %v failed. Want: %v Got: %v", i, c.want, got)
				break
			}
		}
	}
}
