package main

import (
	"testing"
)

func TestGramToCup(t *testing.T) {
	convTable["Quinoa"] = 0.555
	i := Ingr{
		Amount: 1,
		Unit:   cup,
		Item:   "Quinoa",
	}
	i.uoms()
	t.Log(i.AltUnits)
	// want := fmt.Sprintf("0.2775 %v", gram)
	// if got := i.uom(); got != want {
	// 	t.Errorf("Error. Want: %v. Got: %v", want, got)
	// }
}

func TestToTitle(t *testing.T) {
	cases := []struct {
		input, want string
	}{
		{"Test", "Test"},
		{"test", "Test"},
		{"TEST", "Test"},
		{"", ""},
		{"T", "T"},
		{"t", "T"},
	}
	for i, c := range cases {
		got := toTitle(c.input)
		if got != c.want {
			t.Errorf("Case %v failed. Want: %v, Got: %v", i, c.want, got)
		}
	}
}
