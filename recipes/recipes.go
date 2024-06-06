package gocookbook

import (
	"errors"
	"sort"
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

const idSteps = 10 // idSteps is the increment that is used for each new Recipe ID. E.g. if idSteps is 10, then IDs will be 10, 20, 30. If it is 12, then: 12, 24, 36.

var (
	errorUnknownRecipe = errors.New("recipe not found") // Not Found Error.
)

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
func (cb *Cookbook) Add(r Recipe) int {
	r.Id = newRecipeId(*cb)
	*cb = append(*cb, r)
	return r.Id
}

// Find takes an ID, finds the Find in the Cookbook with that ID and returns the Find.
func (cb Cookbook) Find(id int) (Recipe, error) {
	rp, err := findRecipe(cb, id)
	return *rp, err
}

// findRecipe takes a Cookbook of recipes and an id. It looks up the recipe with that id and returns the recipe.
// If the recipe does not exist, it returns an empty Recipe and an error.
func findRecipe(cb Cookbook, id int) (*Recipe, error) {
	for i := range cb {
		if cb[i].Id == id {
			return &cb[i], nil
		}
	}
	return &Recipe{}, errorUnknownRecipe
}

// Update takes a recipe ID and all recipe parameters that can be updated. It finds the recipe for that ID and updates the Recipe.
func (cb *Cookbook) Update(id int, rNew Recipe) error {
	var err error
	r, err := findRecipe(*cb, id)
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

// adjustRcp adjusts the amount of all Ingredients in the Recipe r to the desired portions and returns the adjusted Recipe.
func adjustRcp(r Recipe, portions float64) Recipe {
	newIngrs := make([]Ingredient, len(r.Ingrs))
	copy(newIngrs, r.Ingrs)
	newRcp := Recipe{
		Id:         r.Id,
		Name:       r.Name,
		Ingrs:      newIngrs,
		Steps:      r.Steps,
		Tags:       r.Tags,
		Portions:   portions,
		Dur:        r.Dur,
		Source:     r.Source,
		SourceLink: r.SourceLink,
	}
	x := round(portions / r.Portions)
	for i, v := range newRcp.Ingrs {
		newRcp.Ingrs[i].Amount = round(v.Amount * x)
	}
	return newRcp
}

// FindIngredient takes a string and search in the Cookbook if any Recipe contains an Ingredient 
// where the Item contains the string and returns a Cookbook with alle those recipes.
func (cb Cookbook) FindIngredient(item string) Cookbook{
	result := findIngr(cb, item)
	return result
}

// FindIngr takes a slice of recipes and an item. It returns all recipes that
// have an ingredient that (partially) matches the item and/or recipe name.
func findIngr(rcps Cookbook, item string) Cookbook {
	item = strings.ToLower(item)
	var output Cookbook
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

// Remove takes an Recipe id. The recipe that matches the id is removed
// from the Cookbook and an error is returned (nil if succesful).
func (cb *Cookbook) Remove(id int) error {
	cbOld := *cb
	for i, rcp := range cbOld {
		if rcp.Id == id {
			cbOld[i] = cbOld[len(cbOld)-1] // replace Recipe on index i with the last entry in the Cookbook
			cbNew := cbOld[:len(cbOld)-1] // Cookbook where the Recipe is removed
			sort.SliceStable(cbNew, func(i, j int) bool { return cbNew[i].Name < cbNew[j].Name })
			*cb = cbNew
			return nil
		}
	}
	return errorUnknownRecipe
}
