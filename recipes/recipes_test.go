package gocookbook

import (
	"reflect"
	"testing"
	"time"
	"errors"
)

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
		ckb.Add(want.Name, want.Ingrs, want.Steps, want.Tags, want.Portions, want.Dur, want.Notes, want.Source, want.SourceLink, want.Createdby)
		want.Created = ckb[0].Created
		want.Updated = ckb[0].Updated
		if !reflect.DeepEqual(ckb[0], want) {
			t.Errorf("\nGot: '%+v'\n, Want: '%+v'", ckb[0], want)
		}
	})
	t.Run("update recipe", func(t *testing.T) {
		newName := "test 2"
		newUpdatedBy := "Tester X"
		err := ckb.Update(ckb[0].Id, newName, ckb[0].Ingrs,ckb[0].Steps, ckb[0].Tags, ckb[0].Portions, ckb[0].Dur, ckb[0].Notes, ckb[0].Source, ckb[0].SourceLink, newUpdatedBy)
		switch {
		case errors.Is(err, errorUnknownRecipe):
			t.Errorf("Recipe ID '%v' does not exist: %v", ckb[0].Id, err)
		case err != nil:
			t.Errorf("Unknown error during update of ID %v: %v", err, ckb[0].Id)
		case ckb[0].Name != newName:
			t.Errorf("Want: '%v', Got '%v'", newName, ckb[0].Name)
		case ckb[0].Updatedby != newUpdatedBy:
			t.Errorf("Want: '%v', Got '%v'", newUpdatedBy, ckb[0].Updatedby)
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
