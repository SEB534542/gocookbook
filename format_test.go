package main

import (
	"testing"
)

func TestFormat(t *testing.T) {
	s := `
1 tablespoon extra-virgin olive oil

1 cup thinly sliced celery

1 cup chopped carrots

½ cup chopped onions

8 ounces button mushrooms, sliced

¼ cup all-purpose flour

½ teaspoon ground pepper

½ teaspoon salt

4 cups low-sodium vegetable broth

2 cups cooked wild rice

½ cup heavy cream

2 teaspoons lemon juice

2 tablespoons chopped fresh parsley`
	xi := textToIngrds(s)
	if len(xi) != 8 || xi[0].Amount != 1 {
		t.Error("Ingredients not correctly converted!")
	}
}
