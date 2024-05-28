package gocookbook

import (
	"golang.org/x/text/unicode/norm"
	"strconv"
	"strings"
)

var unitsconv = map[string]func(float64) (Unit, float64){
	"gram":        func(f float64) (Unit, float64) { return gram, f },
	"grams":       func(f float64) (Unit, float64) { return gram, f },
	"gr":          func(f float64) (Unit, float64) { return gram, f },
	"gr.":         func(f float64) (Unit, float64) { return gram, f },
	"g":           func(f float64) (Unit, float64) { return gram, f },
	"g.":          func(f float64) (Unit, float64) { return gram, f },
	"kilogram":    func(f float64) (Unit, float64) { return gram, f * 1000 },
	"kg":          func(f float64) (Unit, float64) { return gram, f * 1000 },
	"kg.":         func(f float64) (Unit, float64) { return gram, f * 1000 },
	"cup":         func(f float64) (Unit, float64) { return cup, f },
	"cups":        func(f float64) (Unit, float64) { return cup, f },
	"ml":          func(f float64) (Unit, float64) { return ml, f },
	"milliliter":  func(f float64) (Unit, float64) { return ml, f },
	"liter":       func(f float64) (Unit, float64) { return ml, f * 1000 },
	"tbsp":        func(f float64) (Unit, float64) { return tbsp, f },
	"tablespoon":  func(f float64) (Unit, float64) { return tbsp, f },
	"tablespoons": func(f float64) (Unit, float64) { return tbsp, f },
	"el":          func(f float64) (Unit, float64) { return tbsp, f },
	"el.":         func(f float64) (Unit, float64) { return tbsp, f },
	"eetlepel":    func(f float64) (Unit, float64) { return tbsp, f },
	"eetlepels":   func(f float64) (Unit, float64) { return tbsp, f },
	"tsp":         func(f float64) (Unit, float64) { return tsp, f },
	"teaspoon":    func(f float64) (Unit, float64) { return tsp, f },
	"teaspoons":   func(f float64) (Unit, float64) { return tsp, f },
	"tl":          func(f float64) (Unit, float64) { return tsp, f },
	"tl.":         func(f float64) (Unit, float64) { return tsp, f },
	"theelepel":   func(f float64) (Unit, float64) { return tsp, f },
	"theelepels":  func(f float64) (Unit, float64) { return tsp, f },
	"stuk":        func(f float64) (Unit, float64) { return pcs, f },
	"stuks":       func(f float64) (Unit, float64) { return pcs, f },
	"pieces":      func(f float64) (Unit, float64) { return pcs, f },
	"pcs":         func(f float64) (Unit, float64) { return pcs, f },
}

/*
	textToLines takes a string, splits the string into a slice for each new line

and removes all non text characters and empty lines. It returns the slice.
*/
func textToLines(s string) []string {
	s = norm.NFC.String(s)

	// Change CR into LR to ensure all 'enters' are split into lines
	s = strings.ReplaceAll(s, "\r", "\n")

	// Change No-Break Spaces into normal spaces
	s = strings.ReplaceAll(s, "\u202f", " ")
	s = strings.ReplaceAll(s, "\u00a0", " ")

	// Split string into lines
	lines := strings.Split(s, "\n")

	// Remove empty lines
	newLines := []string{}
	for _, line := range lines {
		if line != "" {
			newLines = append(newLines, line)
		}
	}
	return newLines
}

/* textToIngrds takes a string and returns a slice of ingredients in the text.*/
func textToIngrds(s string) []Ingredient {
	lines := textToLines(s)
	xi := []Ingredient{}
	// Convert each line to an ingredient
	for _, line := range lines {
		var i Ingredient
		xs := strings.Split(line, " ")
		// Parse each element of the line to a float to find the amount
		for j, s := range xs {
			amount, err := strconv.ParseFloat(s, 64)
			if strings.Index(s, string(uint8(189))) != -1 {
				amount = 0.5
				err = nil
			}
			if strings.Index(s, string(uint8(188))) != -1 {
				amount = 0.25
				err = nil
			}
			if err == nil {
				// Check if a unit is included directly behind the float
				offset := 0
				unit := pcs
				_, ok := unitsconv[xs[j+1]]
				if ok {
					unit, amount = unitsconv[xs[j+1]](amount)
					offset += len(xs[j+1])
				}
				i = Ingredient{
					Amount: amount,
					Unit:   unit,
					Item:   strings.Trim(line[strings.Index(line, s)+len(s)+offset+1:], " "), // assuming item is directly after amount in the text
					Notes:  strings.Trim(line[:strings.Index(line, s)], " "),                 // assuming an text before the float is additional notes
				}
				break
			}
		}
		if i.Amount == 0 && i.Unit == "" && i.Item == "" {
			// No amount found
			i.Item = line
		}
		xi = append(xi, i)
	}
	return xi
}
