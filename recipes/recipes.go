package gocookbook

import (
	"fmt"
	"sort"
	//"strconv"
	"errors"
	"strings"
	"time"
)

type Cookbook []Recipe

// Recipe represents an actual recipe for cooking.
type Recipe struct {
	Id         int           // Internal reference number for a recipe
	Name       string        // Name of recipe.
	Ingrs      []Ingrd       // Slice containing all ingredients.
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

// Ingr represents an ingredient for a recipe.
type Ingrd struct {
	Amount   float64 // Amount of units.
	Unit     string  // Unit of Measurement (UOM), e.g. grams etc.
	Item     string  // Item itself, e.g. a banana.
	Notes    string  // Instruction for preparation, e.g. cooked.
	AltUnits string  // Alternative UOM and the required amount for that unit.
}

var (
	errorUnknownRecipe = errors.New("recipe not found") // Not Found Error
)

// Newcookbook creates a new empty cookbook and returns it.
func NewCookbook() Cookbook {
	return Cookbook{}
}

// Add takes all parameters for a Recipe, creates a new Recipe with that information and adds it to the Cookbook.
func (ckb *Cookbook) Add(name string, ingrs []Ingrd, steps, tags []string, portions float64, dur time.Duration, notes, source, sourceLink, createdby string) {
	rcp := Recipe{
		Id:         newRcpId(*ckb),
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
	*ckb = append(*ckb, rcp)
}

// Recipe takes an ID, finds the Recipe in the Cookbook with that ID and returns the Recipe.
func (ckb Cookbook) Recipe(id int) (Recipe, error) {
	rp, err := findRecipe(ckb, id)
	return *rp, err
}

// Update takes a recipe ID and all recipe parameters that can be updated. It finds the recipe for that ID and updates the Recipe.
func (ckb *Cookbook) Update (id int, name string, ingrs []Ingrd, steps, tags []string, portions float64, dur time.Duration, notes, source, sourceLink, updatedby string) error {
	var err error
	rcp, err := findRecipe(*ckb, id)
	if err != nil {
		return err
	}
	rcp.Update(name, ingrs, steps, tags, portions, dur, notes, source, sourceLink, updatedby)
	return err
}

//newRcpId takes a Cookbook, looks up the highest recipe Id and returns a new recipe Id
func newRcpId(ckb Cookbook) int {
	const idSteps = 10 // increment that is used for each new ID. E.g. if idSteps is 10, then IDs will be 10, 20, 30. If it is 12, then: 12, 24, 36
	var maxId int
	for _, v := range ckb {
		if v.Id > maxId {
			maxId = v.Id
		}
	}
	return maxId + idSteps
}

// Update takes all parameters that can be updated and updates the Recipe pointer.
func (r *Recipe) Update(name string, ingrs []Ingrd, steps, tags []string, portions float64, dur time.Duration, notes, source, sourceLink, updatedby string){
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

// findRecipe takes a Cookbook of recipes and an id. It looks up the recipe with that id and returns the recipe. If the recipe does not exist, it returns an empty Recipe and an error.
func findRecipe(ckb Cookbook, id int) (*Recipe, error) {
	for i := range ckb {
		if ckb[i].Id == id {
			return &ckb[i], nil
		}
	}
	return &Recipe{}, errorUnknownRecipe
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

//Print returns the ingredient with all available information (depending on the type of ingredient) as a string.
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
