package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Recipe represents an actual recipe for cooking.
type Recipe struct {
	Id         int           // Internal reference number for a recipe
	Name       string        // Name of recipe.
	Ingrs      []Ingrd       // Slice containing all ingredients.
	Steps      []string      // Steps for cooking the recipe.
	Tags       []string      // Tags for a recipe.
	Persons    int           // Default number of persons for which this recipe is made.
	Dur        time.Duration // Cooking time
	Notes      string        // Notes and/or description on recipes.
	Source     string        // Source of the recipe.
	SourceLink string        // Hyperlink to the source.
	CreatedBy  string        // User that created the recipe.
	Created    time.Time     // Datetime when created.
	UpdatedBy  string        // User that last updated the recipe.
	Updated    time.Time     // Datetime when last updated.

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
		Id:         rcp.Id,
		Name:       rcp.Name,
		Ingrs:      newIngrs,
		Steps:      rcp.Steps,
		Tags:       rcp.Tags,
		Persons:    newP,
		Dur:        rcp.Dur,
		Source:     rcp.Source,
		SourceLink: rcp.SourceLink,
	}
	x := float64(newP) / float64(rcp.Persons)
	for i, v := range newRcp.Ingrs {
		newRcp.Ingrs[i].Amount = round(v.Amount * x)
	}
	return newRcp
}

func (i Ingrd) Print() string {
	i.uoms()
	var s string
	if i.Unit == pcs {
		s = fmt.Sprintf("%v %v", i.Amount, i.Item)
	} else {
		s = fmt.Sprintf("%v %v %v", i.Amount, i.Unit, i.Item)
	}
	if i.Notes != "" {
		s = fmt.Sprintf("%v, %v", s, strings.ToLower(i.Notes))
	}
	if i.AltUnits != "" {
		return fmt.Sprintf("%v (%v)", s, i.AltUnits)
	}
	return s
}

/* FindIngr takes a slice of recipes and an item. It returns all recipes that
have an ingredient that (partially) matches the item and/or recipe name.*/
func findIngr(rcps []Recipe, item string) []Recipe {
	item = strings.ToLower(item)
	var output []Recipe
	for _, rcp := range rcps {
		if strings.Contains(strings.ToLower(rcp.Name), item) {
			output = append(output, rcp)
			continue
		}
		for _, ingrd := range rcp.Ingrs {
			if strings.Contains(strings.ToLower(ingrd.Item), item) {
				output = append(output, rcp)
				break
			}
		}
	}
	return output
}

/* RemoveRecipe takes a slice of recipes and an id. The recipe that matches
the id is removed from the slice and slice is returned. */
func removeRecipe(rcps []Recipe, id int) []Recipe {
	var i int
	var rcp Recipe
	var b bool
	for i, rcp = range rcps {
		if rcp.Id == id {
			b = true
			break
		}
	}
	if b {
		rcps[i] = rcps[len(rcps)-1]
		rcpsNew := rcps[:len(rcps)-1]
		sort.Slice(rcpsNew, func(i, j int) bool { return rcpsNew[i].Name < rcpsNew[j].Name })
		return rcpsNew
	}
	return rcps
}
