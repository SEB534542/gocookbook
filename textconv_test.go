package main

import (
	"testing"
)

func TestSpecialChars(t *testing.T) {
	s := `300 g couscous 

1 rode ui 

2 teentjes knoflook 

100 g zwarte olijven 

100 g zongedroogde tomaten 

2 venkelknollen 

1 citroen 

olijfolie om in te bakken 

300 g spinazie 

225 g halloumi 

100 g blanke amandelen `
	//lines := textToLines(s)
	lines := textToIngrds(s)
	for _, v := range lines {
		t.Log(v)
	}
}

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

2 tablespoons chopped fresh parsley`
	xi := textToIngrds(s)
	// for _, v := range xi {
	// 	t.Logf("'%v', '%v', '%v', '%v'", v.Amount, v.Unit, v.Item, v.Notes)
	// }
	if len(xi) != 12 || xi[0].Amount != 1 {
		t.Error("Ingredients not correctly converted!")
	}
}
