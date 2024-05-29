package gocookbook

import (
	"strconv"
	"strings"

	"golang.org/x/text/unicode/norm"
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

// TextToIngrds takes a string containing multiple lines of ingredients and returns a slice of ingredients in the text.
func TextToIngrds(s string) []Ingredient {
	lines := textToLines(s)
	xi := make([]Ingredient, len(lines))
	// Convert each line to an ingredient
	for i, line := range lines {
		var in Ingredient
		xs := strings.Split(line, " ")
		// Parse each element of the line to a float to find the amount
		for j, s := range xs {
			amount, err := strconv.ParseFloat(s, 64)
			// check if it is a character for a fractal value
			switch {
			case strings.Contains(s, string(uint8(189))): // character for '1/2'
				amount, err = 0.5, nil
			case strings.Contains(s, string(uint8(188))): // character for '1/4'
				amount, err = 0.25, nil
			case strings.Contains(s, string(uint8(190))): // character '1/3'
				amount, err = 0.33, nil
			}
			if err == nil {
				// Check if a unit is included directly behind the float
				offset := 0
				unit := pcs // default unit if not identified
				_, ok := unitsconv[xs[j+1]]
				if ok {
					unit, amount = unitsconv[xs[j+1]](amount)
					offset += len(xs[j+1])
				}
				item := strings.Trim(line[strings.Index(line, s)+len(s)+offset+1:], " ") // assuming item is directly after amount in the text
				notes := strings.Trim(line[:strings.Index(line, s)], " ")  // assuming an text before the float is additional notes
				
				// if notes is empty, check for a comma in the item and use the remainder as a note
				if notes == "" {
					if x := strings.Index(item, ","); x != -1 {
						if x+1 <= len(item) {
							notes = strings.Trim(item[x+1:], " ")
						}
						item = strings.Trim(item[:x], " ")
					}
				}
				in = NewIngredient(amount, unit, item, notes)
				break
			}
		}
		if in.Amount == 0 && in.Unit == "" && in.Item == "" {
			// No amount found
			in.Item = line
		}
		xi[i] = in
	}
	return xi
}

// textToLines takes a string, splits the string into a slice for each new line and removes all non text characters and empty lines. It returns the slice.
func textToLines(s string) []string {
	s = norm.NFC.String(s)

	// Change CR into LR to ensure all 'enters' are split into lines
	s = strings.ReplaceAll(s, "\r", "\n")

	// Change No-Break Spaces into normal spaces
	s = strings.ReplaceAll(s, "\u202f", " ")
	s = strings.ReplaceAll(s, "\u00a0", " ")

	// Remove tabs
	s = strings.ReplaceAll(s, "\t", "")

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
