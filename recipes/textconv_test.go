package gocookbook

import (
	"testing"
)

func TestTextToIngrds(t *testing.T) {
	want := []Ingredient{
		{
			Amount:   1,
			Unit:     tbsp,
			Item:     "extra-virgin olive oil",
			Notes:    "",
			AltUnits: "14.8 ml / 0.1 cup",
		},
		{
			Amount:   1,
			Unit:     cup,
			Item:     "thinly sliced celery",
			Notes:    "",
			AltUnits: "236.6 ml",
		},
		{
			Amount:   1,
			Unit:     cup,
			Item:     "carrots",
			Notes:    "chopped",
			AltUnits: "236.6 ml",
		},
		{
			Amount:   0.5,
			Unit:     cup,
			Item:     "chopped onions",
			Notes:    "",
			AltUnits: "118.3 ml",
		},
	}
	s := "\n1 tablespoon extra-virgin olive oil\n\n1 cup thinly sliced celery\n\n1 cup carrots, chopped\n\n½ cup chopped onions\n\n8 ounces button mushrooms, sliced\n\n¼ cup all-purpose flour\n\n½ teaspoon ground pepper\n\n½ teaspoon salt\n\n4 cups low-sodium vegetable broth\n\n2 cups cooked wild rice\n\n½ cup heavy cream\n\n2 tablespoons chopped fresh parsley"
	got := TextToIngrds(s)
	if len(got) != 12 {
		t.Errorf("not all ingredients are converted. Want: %v, Got: %v", 12, len(got))
	}
	for i := range want {
		if b, fields := AssertEqualIngrd(want[i], got[i]); !b {
			t.Errorf("fields %v are not equal for item %v:\nGot:\t'%+v'\nWant:\t'%+v'", fields, i, got[i], want[i])
		}
	}
}
