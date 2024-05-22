package gocookbook

import (
	"testing"
	"reflect"
)

func TestNewCookbook(t *testing.T){
	ckb := NewCookbook()
	t.Run("new cookbook", func(t *testing.T){
		want := Cookbook{}
		if !reflect.DeepEqual(ckb, want) {
			t.Errorf("Got: '%v', Want: '%v'", ckb, want)
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
