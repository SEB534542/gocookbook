package main

import (
	"fmt"
)

// https://simoneskitchen.nl/wprm_print/recipe/43253
// Recipemaker: https://simoneskitchen.nl/wp-content/plugins/wp-recipe-maker/dist/print.js?ver=8.3.0

// Recipe represents an actual recipe for cooking.
type Recipe struct {
	Id      int      // Internal reference number for a recipe
	Name    string   // Name of recipe.
	Ingrs   []Ingr   // Slice containing all ingredients.
	Steps   []string // Steps for cooking the recipe.
	Persons int      // Default number of persons for which this recipe is made.
	Notes   string   // Notes and/or description on recipes.
	Source  string   // Source of the recipe.
}

// TODO: implement logic for tags

// Ingr represents an ingredient for a recipe.
type Ingr struct {
	Amount float64 // Amount of units.
	Unit   string  // Unit of Measurement (UOM), e.g. grams etc. TODO: make uom a tye?
	Item   string  // Item itself, e.g. a banana.
	Notes  string  // Instruction for preparation, e.g. cooked.
	Alts   string  // Alternative UOM and the required amount for that unit.
}

var (
	errorUnknownRecipe = fmt.Errorf("Recipe not found.")
)

var fnameRcps = "recipes.json"

var rcps []Recipe

func main() {
	err := readJSON(&rcps, fnameRcps)
	if err != nil {
		fmt.Println(err)
	}
	startServer(8081)
}

func findRecipe(rcps []Recipe, id int) (Recipe, error) {
	for _, rcp := range rcps {
		if rcp.Id == id {
			return rcp, nil
		}
	}
	return Recipe{}, errorUnknownRecipe
}

func newId(rcps []Recipe) int {
	var maxId int
	for _, v := range rcps {
		if v.Id > maxId {
			maxId = v.Id
		}
	}
	return maxId + 10
}

func findRecipeP(rcps []Recipe, id int) (*Recipe, error) {
	for i, _ := range rcps {
		if rcps[i].Id == id {
			return &rcps[i], nil
		}
	}
	return &Recipe{}, errorUnknownRecipe
}

// updateRcp adjusts Ingrs in the recipe r to n persons and returns the new recipe.
func adjustRcp(rcp Recipe, newP int) Recipe {
	newIngrs := make([]Ingr, len(rcp.Ingrs))
	copy(newIngrs, rcp.Ingrs)
	newRcp := Recipe{
		Id:      rcp.Id,
		Name:    rcp.Name,
		Ingrs:   newIngrs,
		Steps:   rcp.Steps,
		Persons: newP,
		Source:  rcp.Source,
	}
	x := float64(newP) / float64(rcp.Persons)
	for i, _ := range newRcp.Ingrs {
		newRcp.Ingrs[i].Amount *= x
	}
	return newRcp
}
