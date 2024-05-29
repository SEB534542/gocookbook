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
	ckb := NewCookbook()
	t.Run("new cookbook", func(t *testing.T) {
		want := Cookbook{}
		if !reflect.DeepEqual(ckb, want) {
			t.Errorf("Got: '%v', Want: '%v'", ckb, want)
		}
	})
	t.Run("Add test recipes", func(t *testing.T) {
		want := NewRecipe("Test1", []Ingredient{}, []string{}, []string{}, 4, 0, "", "", "", "Tester1")
		want.Id = 0 + idSteps
		ckb.Add(want)
		want.Created, want.Updated = ckb[0].Created, ckb[0].Updated
		if !reflect.DeepEqual(ckb[0], want) {
			t.Errorf("\nGot: '%+v'\n, Want: '%+v'", ckb[0], want)
		}

		want2 := want
		want2.Name = "Test2"
		want2.Id = idSteps * 2
		ckb.Add(want2)
		want2.Created, want2.Updated = ckb[1].Created, ckb[1].Updated
		if !reflect.DeepEqual(ckb[1], want2) {
			t.Errorf("\nGot: '%+v'\n, Want: '%+v'", ckb[1], want2)
		}
	})
	t.Run("update existing recipe", func(t *testing.T) {
		want := ckb[0]
		want.Update("test1 v2", []Ingredient{}, []string{}, []string{}, 5, 0, "", "", "", "Tester 2")
		err := ckb.Update(idSteps, want)
		want.Updated = ckb[0].Updated
		switch {
		case errors.Is(err, errorUnknownRecipe):
			t.Errorf("Recipe ID '%v' does not exist: %v", ckb[0].Id, err)
		case err != nil:
			t.Errorf("Unknown error during update of ID %v: %v", err, ckb[0].Id)
		case ckb[0].Name != want.Name:
			t.Errorf("Want: '%v', Got '%v'", want.Name, ckb[0].Name)
		case ckb[0].Updatedby != want.Updatedby:
			t.Errorf("Want: '%v', Got '%v'", want.Updatedby, ckb[0].Updatedby)
		}
	})
	t.Run("create recipe with ingredients from string", func(t *testing.T) {
		s := "\n\t\t1 tablespoon extra-virgin olive oil\n\t\t\n\t\t1 cup thinly sliced celery\n\t\t\n\t\t1 cup chopped carrots\n\t\t\n\t\t½ cup chopped onions\n\t\t\n\t\t8 ounces button mushrooms, sliced\n\t\t\n\t\t¼ cup all-purpose flour\n\t\t\n\t\t½ teaspoon ground pepper\n\t\t\n\t\t½ teaspoon salt\n\t\t\n\t\t4 cups low-sodium vegetable broth\n\t\t\n\t\t2 cups cooked wild rice\n\t\t\n\t\t½ cup heavy cream\n\t\t\n\t\t2 tablespoons chopped fresh parsley"
		id := ckb.Add(NewRecipe(
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
		_, err := ckb.Recipe(id)
		switch {
		case errors.Is(err, errorUnknownRecipe):
			t.Errorf("Recipe ID '%v' does not exist: %v", id, err)
		case err != nil:
			t.Errorf("Unknown error retrieving recipe ID %v: %v", err, id)
		}
		//if AssertEqualIngrd()
		//t.Logf("%+v", r.Ingrs[0])
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
	if len(removeRecipe(rcps, 2)) != 2 {
		t.Error("Recipe not removed")
	}
}
