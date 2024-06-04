package gocookbook

import (
	"testing"
)

func TestNewIngr(t *testing.T) {
	want := Ingredient{
		Amount:   50,
		Unit:     ml,
		Item:     "Banaan",
		Notes:    "in plakjes",
		AltUnits: "0.2 cup",
	}
	in := NewIngredient(want.Amount, want.Unit, want.Item, want.Notes)

	if b, fields := AssertEqualIngrd(want, in); !b {
		t.Errorf("fields %v are not equal\nGot:\t'%+v'\nWant:\t'%+v'", fields, in, want)
	}
}

func AssertEqualIngrd(x, y Ingredient) (bool, []string) {
	var b bool
	var fields []string
	if x.Amount != y.Amount {
		fields = append(fields, "Amount")
	}
	if x.Unit != y.Unit {
		fields = append(fields, "Unit")
	}
	if x.Item != y.Item {
		fields = append(fields, "Item")
	}
	if x.Notes != y.Notes {
		fields = append(fields, "Notes")
	}
	if x.AltUnits != y.AltUnits {
		fields = append(fields, "AltUnits")
	}
	if len(fields) == 0 {
		b = true
	}
	return b, fields
}

func TestPrint(t *testing.T){
	in := NewIngredient(50, ml, "Banaan", "in plakjes")
	want := "50 ml Banaan, in plakjes (0.2 cup)"
	if got := in.Print(); want != got {
		t.Errorf("Want: '%v', Got: '%v'", want, got)
	}
}