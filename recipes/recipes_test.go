package gocookbook

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestRecipe(t *testing.T) {
	want := Recipe{
		Id:         0,
		Name:       "Test1",
		Ingrs:      []Ingredient{},
		Steps:      []string{},
		Tags:       []string{},
		Portions:   4,
		Dur:        0,
		Notes:      "",
		Source:     "",
		SourceLink: "",
		Createdby:  "Tester1",
		Created:    time.Time{},
		Updatedby:  "Tester1",
		Updated:    time.Time{},
	}
	r := NewRecipe(want.Name, want.Ingrs, want.Steps, want.Tags, want.Portions, want.Dur, want.Notes, want.Source, want.SourceLink, want.Createdby)
	want.Created = r.Created
	want.Updated = r.Updated
	if !reflect.DeepEqual(r, want) {
		t.Errorf("\nGot: '%+v'\n, Want: '%+v'", r, want)
	}

	t.Run("update recipe", func(t *testing.T) {
		want.Name = "Test2"
		want.Updatedby = "Tester2"
		r.Update(want.Name, want.Ingrs, want.Steps, want.Tags, want.Portions, want.Dur, want.Notes, want.Source, want.SourceLink, want.Updatedby)
		if want.Updated == r.Updated {
			t.Errorf("updated time is not updated. Want: %v, Got: %v", want.Updated, r.Updated)
		}
		want.Updated = r.Updated
		if !reflect.DeepEqual(r, want) {
			t.Errorf("\nGot: '%+v'\n, Want: '%+v'", r, want)
		}
	})
}

func TestNewCookbook(t *testing.T) {
	cb := NewCookbook()
	t.Run("new cookbook", func(t *testing.T) {
		want := Cookbook{}
		if !reflect.DeepEqual(cb, want) {
			t.Errorf("Got: '%v', Want: '%v'", cb, want)
		}
	})
	t.Run("Add test recipes", func(t *testing.T) {
		want := NewRecipe("Test1", []Ingredient{}, []string{}, []string{}, 4, 0, "", "", "", "Tester1")
		want.Id = 0 + idSteps
		cb.Add(want)
		want.Created, want.Updated = cb[0].Created, cb[0].Updated
		if !reflect.DeepEqual(cb[0], want) {
			t.Errorf("\nGot: '%+v'\n, Want: '%+v'", cb[0], want)
		}

		want2 := want
		want2.Name = "Test2"
		want2.Id = idSteps * 2
		cb.Add(want2)
		want2.Created, want2.Updated = cb[1].Created, cb[1].Updated
		if !reflect.DeepEqual(cb[1], want2) {
			t.Errorf("\nGot: '%+v'\n, Want: '%+v'", cb[1], want2)
		}
	})
	t.Run("update existing recipe", func(t *testing.T) {
		want := cb[0]
		want.Update("test1 v2", []Ingredient{}, []string{}, []string{}, 5, 0, "", "", "", "Tester 2")
		err := cb.Update(idSteps, want)
		want.Updated = cb[0].Updated
		switch {
		case errors.Is(err, errorUnknownRecipe):
			t.Errorf("Recipe ID '%v' does not exist: %v", cb[0].Id, err)
		case err != nil:
			t.Errorf("Unknown error during update of ID %v: %v", err, cb[0].Id)
		case cb[0].Name != want.Name:
			t.Errorf("Want: '%v', Got '%v'", want.Name, cb[0].Name)
		case cb[0].Updatedby != want.Updatedby:
			t.Errorf("Want: '%v', Got '%v'", want.Updatedby, cb[0].Updatedby)
		}
	})
	t.Run("create recipe with ingredients from string", func(t *testing.T) {
		s := "\n\t\t1 tablespoon extra-virgin olive oil\n\t\t\n\t\t1 cup thinly sliced celery\n\t\t\n\t\t1 cup chopped carrots\n\t\t\n\t\t½ cup chopped onions\n\t\t\n\t\t8 ounces button mushrooms, sliced\n\t\t\n\t\t¼ cup all-purpose flour\n\t\t\n\t\t½ teaspoon ground pepper\n\t\t\n\t\t½ teaspoon salt\n\t\t\n\t\t4 cups low-sodium vegetable broth\n\t\t\n\t\t2 cups cooked wild rice\n\t\t\n\t\t½ cup heavy cream\n\t\t\n\t\t2 tablespoons chopped fresh parsley"
		id := cb.Add(NewRecipe(
			"test3",
			TextToIngrds(s),
			[]string{},
			[]string{},
			4,
			0,
			"",
			"",
			"",
			"Tester 1",
		))
		_, err := cb.Find(id)
		switch {
		case errors.Is(err, errorUnknownRecipe):
			t.Errorf("Recipe ID '%v' does not exist: %v", id, err)
		case err != nil:
			t.Errorf("Unknown error retrieving recipe ID %v: %v", id, err)
		}
		//if AssertEqualIngrd()
		//t.Logf("%+v", r.Ingrs[0])

		t.Run("find ingredient", func(t *testing.T) {
			want := "extra-virgin olive oil"
			result := cb.FindIngredient(want)
			if got := result[0].Ingrs[0].Item; got != want {
				t.Errorf("Want: %v, Got: %v", want, got)
			}
		})
	})
}

func TestFindIngr(t *testing.T) {
	rcps := []Recipe{
		{
			Id:   1,
			Name: "Test1",
			Ingrs: []Ingredient{
				{Item: "Paddestoelenboullion"},
			},
		},
		{
			Id:   2,
			Name: "Test2",
			Ingrs: []Ingredient{
				{Item: "Pasta"},
			},
		},
		{
			Id:   3,
			Name: "Test3",
			Ingrs: []Ingredient{
				{Item: "bospaddestoelen"},
			},
		},
	}
	if len(findIngr(rcps, "Paddestoel")) != 2 {
		t.Error("FindIngr no longer works")
	}
}

func TestRemoveRecipe(t *testing.T) {
	cb := NewCookbook()
	id := idSteps
	s := "\n\t\t1 tablespoon extra-virgin olive oil\n\t\t\n\t\t1 cup thinly sliced celery\n\t\t\n\t\t1 cup chopped carrots\n\t\t\n\t\t½ cup chopped onions\n\t\t\n\t\t8 ounces button mushrooms, sliced\n\t\t\n\t\t¼ cup all-purpose flour\n\t\t\n\t\t½ teaspoon ground pepper\n\t\t\n\t\t½ teaspoon salt\n\t\t\n\t\t4 cups low-sodium vegetable broth\n\t\t\n\t\t2 cups cooked wild rice\n\t\t\n\t\t½ cup heavy cream\n\t\t\n\t\t2 tablespoons chopped fresh parsley"
	cb.Add(NewRecipe("Test1", TextToIngrds(s), []string{}, []string{}, 4, 0, "", "", "", "Tester1"))
	err := cb.Remove(id)
	switch {
	case errors.Is(err, errorUnknownRecipe):
		t.Errorf("Unable to delete Recipe %v: %v", id, err)
	case err != nil:
		t.Errorf("Unknown error deleting recipe ID %v: %v", id, err)
	}
	if _, err = cb.Find(id); err == nil {
		t.Errorf("Recipe %v not deleted from Cookbook: %+v", id, cb)
	}
}

func TestAdjustRecipe(t *testing.T) {
	s := "\n\t\t1 tablespoon extra-virgin olive oil\n\t\t\n\t\t1 cup thinly sliced celery\n\t\t\n\t\t1 cup chopped carrots\n\t\t\n\t\t½ cup chopped onions\n\t\t\n\t\t8 ounces button mushrooms, sliced\n\t\t\n\t\t¼ cup all-purpose flour\n\t\t\n\t\t½ teaspoon ground pepper\n\t\t\n\t\t½ teaspoon salt\n\t\t\n\t\t4 cups low-sodium vegetable broth\n\t\t\n\t\t2 cups cooked wild rice\n\t\t\n\t\t½ cup heavy cream\n\t\t\n\t\t2 tablespoons chopped fresh parsley"
	r := NewRecipe(
		"test3",
		TextToIngrds(s),
		[]string{},
		[]string{},
		4,
		0,
		"",
		"",
		"",
		"Tester 1",
	)
	portionsNew := 8.0
	result := adjustRcp(r, portionsNew)
	for i := range result.Ingrs {
		if got, want := result.Ingrs[i].Amount, (r.Ingrs[i].Amount * portionsNew / r.Portions); got != want {
			t.Errorf("Want: %v, Got: %v for %+v ", want, got, result.Ingrs[i])
		}
	}
}
