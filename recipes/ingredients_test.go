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
	i := NewIngredient(want.Amount, string(want.Unit), want.Item, want.Notes)

	if b, fields := AssertEqualIngrd(want, i); !b {
		t.Errorf("fields %v are not equal\nGot:\t'%+v'\nWant:\t'%+v'", fields, i, want)
	}
}

func AssertEqualIngrd(i, j Ingredient) (bool, []string) {
	var b bool
	var fields []string
	if i.Amount != j.Amount {
		fields = append(fields, "Amount")
	}
	if i.Unit != j.Unit {
		fields = append(fields, "Unit")
	}
	if i.Item != j.Item {
		fields = append(fields, "Item")
	}
	if i.Notes != j.Notes {
		fields = append(fields, "Notes")
	}
	if i.AltUnits != j.AltUnits {
		fields = append(fields, "AltUnits")
	}
	if len(fields) == 0 {
		b = true
	}
	return b, fields
}

// func TestAltUnits(t *testing.T){
// 	type test struct {
// 		unit Unit
// 		amount float64
// 	}
// 	xc := []test{
// 		test(gram, 5.0),
// }
// }
