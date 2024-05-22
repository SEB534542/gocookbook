package gocookbook

import (
	"fmt"
	"sort"
	//"strconv"
	"strings"
	"time"
	"errors"
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

var rcps []Recipe // TODO: remove?

func NewCookbook() Cookbook {
	return Cookbook{}
}

// // NewRecipe takes all date required for a recipe and returns the recipe
// func NewRecipe(name string, ingrs []Ingrd, steps, tags []string, portions float64, dur time.Duration, notes, source, sourceLink, createdby string) Recipe {
// 	rcp := Recipe{}
// 	if id := req.PostFormValue("Id"); id != "" {
// 		rcp.Id, _ = strconv.Atoi(id)
// 	}
// 	rcp.Name = strings.Trim(req.PostFormValue("Name"), " ")
// 	rcp.Notes = strings.Trim(req.PostFormValue("Notes"), " ")
// 	rcp.Dur, _ = time.ParseDuration(fmt.Sprintf("%vm", req.PostFormValue("Dur")))
// 	rcp.Portions, _ = strconv.ParseFloat(req.PostFormValue("Portions"), 64)

// 	t := stringToSlice(req.PostFormValue("Tags"))
// 	rcp.Tags = []string{}
// 	for _, v := range t {
// 		v = strings.Trim(v, " ")
// 		if v != "" {
// 			rcp.Tags = append(rcp.Tags, toTitle(v))
// 		}
// 	}
// 	sort.Strings(rcp.Tags)
// 	// Ingredients
// 	rcp.Ingrs = textToIngrds(req.PostFormValue("Ingrds"))
// 	// Steps
// 	rcp.Steps = textToLines(req.PostFormValue("Steps"))
// 	// Store source and hyperlink
// 	rcp.Source = req.PostFormValue("Source")
// 	rcp.SourceLink = req.PostFormValue("SourceLink")
// 	switch {
// 	case rcp.SourceLink == "" && isHyperlink(rcp.Source):
// 		rcp.SourceLink = rcp.Source
// 	case rcp.Source == "" && isHyperlink(rcp.SourceLink):
// 		rcp.Source = rcp.SourceLink
// 	case !isHyperlink(rcp.SourceLink):
// 		rcp.SourceLink = ""
// 	}
// 	/*Store user and datetime. As this func creates a new recipe,
// 	it sets both AddedBy and UpdatedBy to the same user.
// 	In "upper" logic the AddedBy is restored to the original creator,
// 	if it is an update to existing recipe.*/
// 	if un := currentUser(req); un != "" {
// 		rcp.Createdby = un
// 		rcp.Updatedby = un
// 		t := time.Now()
// 		rcp.Created = t
// 		rcp.Updated = t
// 	}
// 	return rcp
// }

/*
	findRecipe takes a slice of recipes and an id. It looks up the recipe with that

id and returns the recipe.
*/
func findRecipe(rcps []Recipe, id int) (Recipe, error) {
	for _, rcp := range rcps {
		if rcp.Id == id {
			return rcp, nil
		}
	}
	return Recipe{}, errorUnknownRecipe
}

/*
	newRcpId takes a slice of Recipes, looks up the highest recipe Id and

returns a new recipe Id.
*/
func newRcpId(rcps []Recipe) int {
	var maxId int
	for _, v := range rcps {
		if v.Id > maxId {
			maxId = v.Id
		}
	}
	return maxId + 10
}

/*
	findRecipeP takes a slice of recipes and an id. It looks up the recipe with that

id and returns a pointer to the recipe.
*/
func findRecipeP(rcps []Recipe, id int) (*Recipe, error) {
	for i, _ := range rcps {
		if rcps[i].Id == id {
			return &rcps[i], nil
		}
	}
	return &Recipe{}, errorUnknownRecipe
}

// updateRcp adjusts Ingrs in the recipe r to n persons and returns the new recipe.
func adjustRcp(rcp Recipe, newP float64) Recipe {
	newIngrs := make([]Ingrd, len(rcp.Ingrs))
	copy(newIngrs, rcp.Ingrs)
	newRcp := Recipe{
		Id:         rcp.Id,
		Name:       rcp.Name,
		Ingrs:      newIngrs,
		Steps:      rcp.Steps,
		Tags:       rcp.Tags,
		Portions:   newP,
		Dur:        rcp.Dur,
		Source:     rcp.Source,
		SourceLink: rcp.SourceLink,
	}
	x := round(newP / rcp.Portions)
	for i, v := range newRcp.Ingrs {
		newRcp.Ingrs[i].Amount = round(v.Amount * x)
	}
	return newRcp
}

/*
	Print returns the ingredient with all available information

(depending on the type of ingredient) as a string.
*/
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
