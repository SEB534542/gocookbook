package gocookbook

import (
	"fmt"
	"sort"

	//"strconv"
	"errors"
	"strings"
	"time"
)

// Cookbook is a slice of Recipes that represents an actual cookbook of recipes for cooking.
type Cookbook []Recipe

// Recipe represents an actual recipe for cooking.
type Recipe struct {
	Id         int           // Internal reference number for a recipe
	Name       string        // Name of recipe.
	Ingrs      []Ingredient  // Slice containing all ingredients.
	Steps      []string      // Steps for cooking the recipe.
	Tags       []string      // Tags for a recipe.
	Portions   float64       // Default number of portions for recipe.
	Dur        time.Duration // Cooking time
	Notes      string        // Notes and/or description on recipes.
	Source     string        // Source of the recipe.
	SourceLink string        // Hyperlink to the source.
	Createdby  string        // User that created the recipe.
	Created    time.Time     // Datetime when created.
	Updatedby  string        // User that last updated the recipe.
	Updated    time.Time     // Datetime when last updated.
}

// Ingrident represents an ingredient for a recipe.
type Ingredient struct {
	Amount   float64 // Amount of units.
	Unit     Unit    // Unit of Measurement (UOM), e.g. grams etc.
	Item     string  // Item itself, e.g. a banana.
	Notes    string  // Instruction for preparation, e.g. cooked.
	AltUnits string  // Alternative UOM and the required amount for that unit.
}

type Unit string // Unit represents a Unit of Measurement.

const idSteps = 10 // idSteps is the increment that is used for each new ID. E.g. if idSteps is 10, then IDs will be 10, 20, 30. If it is 12, then: 12, 24, 36.

var (
	errorUnknownRecipe = errors.New("recipe not found") // Not Found Error.
)

// NewIngredient takes all parameters for creating an Ingredient, validates all parameters and returns it as an Ingredient.
func NewIngredient(amount float64, unit, item, notes string) Ingredient {
	i := Ingredient{
		Amount:   amount,
		Unit:     Unit(unit),
		Item:     item,
		Notes:    notes,
		AltUnits: "",
	}
	i.altUnits() // Add alt units
	return i
}

// NewRecipe takes all parameters for a Recipe, creates a new Recipe and returns it.
func NewRecipe(name string, ingrs []Ingredient, steps, tags []string, portions float64, dur time.Duration, notes, source, sourceLink, createdby string) Recipe {
	return Recipe{
		Id:         0,
		Name:       name,
		Ingrs:      ingrs,
		Steps:      steps,
		Tags:       tags,
		Portions:   portions,
		Dur:        dur,
		Notes:      notes,
		Source:     source,
		SourceLink: sourceLink,
		Createdby:  createdby,
		Created:    time.Now(),
		Updatedby:  createdby,
		Updated:    time.Now(),
	}
}

// Update takes all parameters that can be updated and updates the Recipe pointer.
func (r *Recipe) Update(name string, ingrs []Ingredient, steps, tags []string, portions float64, dur time.Duration, notes, source, sourceLink, updatedby string) {
	r.Name = name
	r.Ingrs = ingrs
	r.Steps = steps
	r.Tags = tags
	r.Portions = portions
	r.Dur = dur
	r.Notes = notes
	r.Source = source
	r.SourceLink = sourceLink
	r.Updatedby = updatedby
	r.Updated = time.Now()
}

// Newcookbook creates a new empty cookbook and returns it.
func NewCookbook() Cookbook {
	return Cookbook{}
}

// Add takes a Recipe and adds it to the Cookbook.
func (ckb *Cookbook) Add(r Recipe) {
	r.Id = newRecipeId(*ckb)
	*ckb = append(*ckb, r)
}

// Recipe takes an ID, finds the Recipe in the Cookbook with that ID and returns the Recipe.
func (ckb Cookbook) Recipe(id int) (Recipe, error) {
	rp, err := findRecipe(ckb, id)
	return *rp, err
}

// findRecipe takes a Cookbook of recipes and an id. It looks up the recipe with that id and returns the recipe. If the recipe does not exist, it returns an empty Recipe and an error.
func findRecipe(ckb Cookbook, id int) (*Recipe, error) {
	for i := range ckb {
		if ckb[i].Id == id {
			return &ckb[i], nil
		}
	}
	return &Recipe{}, errorUnknownRecipe
}

// Update takes a recipe ID and all recipe parameters that can be updated. It finds the recipe for that ID and updates the Recipe.
func (ckb *Cookbook) Update(id int, rNew Recipe) error {
	var err error
	r, err := findRecipe(*ckb, id)
	if err != nil {
		return err
	}
	*r = rNew
	return err
}

// newRecipeId takes a Cookbook, looks up the highest recipe Id and returns a new recipe Id
func newRecipeId(ckb Cookbook) int {
	var maxId int
	for _, v := range ckb {
		if v.Id > maxId {
			maxId = v.Id
		}
	}
	return maxId + idSteps
}

// TODO: remove below?
// // updateRcp adjusts Ingrs in the recipe r to n persons and returns the new recipe.
// func adjustRcp(rcp Recipe, newP float64) Recipe {
// 	newIngrs := make([]Ingrd, len(rcp.Ingrs))
// 	copy(newIngrs, rcp.Ingrs)
// 	newRcp := Recipe{
// 		Id:         rcp.Id,
// 		Name:       rcp.Name,
// 		Ingrs:      newIngrs,
// 		Steps:      rcp.Steps,
// 		Tags:       rcp.Tags,
// 		Portions:   newP,
// 		Dur:        rcp.Dur,
// 		Source:     rcp.Source,
// 		SourceLink: rcp.SourceLink,
// 	}
// 	x := round(newP / rcp.Portions)
// 	for i, v := range newRcp.Ingrs {
// 		newRcp.Ingrs[i].Amount = round(v.Amount * x)
// 	}
// 	return newRcp
// }

// Print returns the ingredient with all available information (depending on the type of ingredient) as a string.
func (i Ingredient) Print() string {
	i.altUnits()
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

/*
	FindIngr takes a slice of recipes and an item. It returns all recipes that

have an ingredient that (partially) matches the item and/or recipe name.
*/
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

/*
	RemoveRecipe takes a slice of recipes and an id. The recipe that matches

the id is removed from the slice and slice is returned.
*/
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
