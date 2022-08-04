package main

import (
	"fmt"
)

// Recipe represents an actual recipe for cooking.
type Recipe struct {
	Id      int      // Internal reference number for a recipe
	Name    string   // Name of recipe.
	Ingrs   []Ingrd  // Slice containing all ingredients.
	Steps   []string // Steps for cooking the recipe.
	Persons int      // Default number of persons for which this recipe is made.
	Notes   string   // Notes and/or description on recipes.
	Source  string   // Source of the recipe.
}

// Ingr represents an ingredient for a recipe.
type Ingrd struct {
	Amount   float64 // Amount of units.
	Unit     string  // Unit of Measurement (UOM), e.g. grams etc.
	Item     string  // Item itself, e.g. a banana.
	Notes    string  // Instruction for preparation, e.g. cooked.
	AltUnits string  // Alternative UOM and the required amount for that unit.
}

var (
	errorUnknownRecipe = fmt.Errorf("Recipe not found.")
)

var rcps []Recipe

// TODO: implement logic for tags

/* findRecipe takes a slice of recipes and an id. It looks up the recipe with that
id and returns the recipe.*/
func findRecipe(rcps []Recipe, id int) (Recipe, error) {
	for _, rcp := range rcps {
		if rcp.Id == id {
			return rcp, nil
		}
	}
	return Recipe{}, errorUnknownRecipe
}

/* newRcpId takes a slice of Recipes, looks up the highest recipe Id and
returns a new recipe Id.*/
func newRcpId(rcps []Recipe) int {
	var maxId int
	for _, v := range rcps {
		if v.Id > maxId {
			maxId = v.Id
		}
	}
	return maxId + 10
}

/* findRecipeP takes a slice of recipes and an id. It looks up the recipe with that
id and returns a pointer to the recipe.*/
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
	newIngrs := make([]Ingrd, len(rcp.Ingrs))
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
