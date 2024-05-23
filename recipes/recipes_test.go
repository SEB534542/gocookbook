package gocookbook

import (
	"reflect"
	"testing"
	"time"
	"errors"
)

func TestNewRecipe(t *testing.T) {
	want := Recipe{
		Id:         0,
		Name:       "Test1",
		Ingrs:      []Ingrd{},
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
	t.Run("Add test recipe", func(t *testing.T) {
		want := Recipe{
			Id:         10,
			Name:       "Test1",
			Ingrs:      []Ingrd{},
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
		ckb.Add(want)
		want.Created = ckb[0].Created
		want.Updated = ckb[0].Updated
		if !reflect.DeepEqual(ckb[0], want) {
			t.Errorf("\nGot: '%+v'\n, Want: '%+v'", ckb[0], want)
		}
	})
	t.Run("update existing recipe", func(t *testing.T) {
		want := ckb[0]
		want.Name = "test 2"
		want.Updatedby = "Tester 2"
		err := ckb.Update(ckb[0].Id, want)
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
}

func TestFindIngr(t *testing.T) {
	rcps := []Recipe{
		{
			Id:   1,
			Name: "Test1",
			Ingrs: []Ingrd{
				{Item: "Paddestoelenboullion"},
			},
		},
		{
			Id:   2,
			Name: "Test2",
			Ingrs: []Ingrd{
				{Item: "Pasta"},
			},
		},
		{
			Id:   3,
			Name: "Test3",
			Ingrs: []Ingrd{
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
