package main

import (
	"testing"
)

func TestFindIngr(t *testing.T) {
	rcps := []Recipe{
		{
			Name: "Test1",
			Ingrs: []Ingrd{
				{Item: "Paddestoelenboullion"},
			},
		},
		{
			Name: "Test2",
			Ingrs: []Ingrd{
				{Item: "Pasta"},
			},
		},
		{
			Name: "Test3",
			Ingrs: []Ingrd{
				{Item: "bospaddestoelen"},
			},
		},
	}
	if len(findIngr(rcps, "Paddestoel")) != 2 {
		t.Error("FindIngr no longer works")
	}
}
