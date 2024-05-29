package gocookbook

import (
	"testing"
)

func TestFormat(t *testing.T) {
	want := []Ingredient{
		Ingredient{
			Amount:   1,
			Unit:     tbsp,
			Item:     "extra-virgin olive oil",
			Notes:    "",
			AltUnits: "",
		},
	}
	s := "\n1 tablespoon extra-virgin olive oil\n\n1 cup thinly sliced celery\n\n1 cup chopped carrots\n\n½ cup chopped onions\n\n8 ounces button mushrooms, sliced\n\n¼ cup all-purpose flour\n\n½ teaspoon ground pepper\n\n½ teaspoon salt\n\n4 cups low-sodium vegetable broth\n\n2 cups cooked wild rice\n\n½ cup heavy cream\n\n2 tablespoons chopped fresh parsley"
	got := TextToIngrds(s)
	if len(got) != 12 {
		t.Errorf("not all ingredients are converted. Want: %v, Got: %v", 12, len(got))
	}	
	for i, _ := range want {
		if b, fields := AssertEqualIngrd(want[i], got[i]); !b {
			t.Errorf("fields %v are not equal\nGot:\t'%+v'\nWant:\t'%+v'", fields, i, want)
		}
	}
	
	
}
