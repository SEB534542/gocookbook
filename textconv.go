package main

import (
	"golang.org/x/text/unicode/norm"
	"strconv"
	"strings"
)

var unitsconv = map[string]func(float64) (string, float64){
	"gram":        func(f float64) (string, float64) { return gram, f },
	"grams":       func(f float64) (string, float64) { return gram, f },
	"gr":          func(f float64) (string, float64) { return gram, f },
	"gr.":         func(f float64) (string, float64) { return gram, f },
	"g":           func(f float64) (string, float64) { return gram, f },
	"g.":          func(f float64) (string, float64) { return gram, f },
	"kilogram":    func(f float64) (string, float64) { return gram, f * 1000 },
	"kg":          func(f float64) (string, float64) { return gram, f * 1000 },
	"kg.":         func(f float64) (string, float64) { return gram, f * 1000 },
	"cup":         func(f float64) (string, float64) { return cup, f },
	"cups":        func(f float64) (string, float64) { return cup, f },
	"ml":          func(f float64) (string, float64) { return ml, f },
	"milliliter":  func(f float64) (string, float64) { return ml, f },
	"liter":       func(f float64) (string, float64) { return ml, f * 1000 },
	"tbsp":        func(f float64) (string, float64) { return tbsp, f },
	"tablespoon":  func(f float64) (string, float64) { return tbsp, f },
	"tablespoons": func(f float64) (string, float64) { return tbsp, f },
	"el":          func(f float64) (string, float64) { return tbsp, f },
	"el.":         func(f float64) (string, float64) { return tbsp, f },
	"eetlepel":    func(f float64) (string, float64) { return tbsp, f },
	"eetlepels":   func(f float64) (string, float64) { return tbsp, f },
	"tsp":         func(f float64) (string, float64) { return tsp, f },
	"teaspoon":    func(f float64) (string, float64) { return tsp, f },
	"teaspoons":   func(f float64) (string, float64) { return tsp, f },
	"tl":          func(f float64) (string, float64) { return tsp, f },
	"tl.":         func(f float64) (string, float64) { return tsp, f },
	"theelepel":   func(f float64) (string, float64) { return tsp, f },
	"theelepels":  func(f float64) (string, float64) { return tsp, f },
	"stuk":        func(f float64) (string, float64) { return pcs, f },
	"stuks":       func(f float64) (string, float64) { return pcs, f },
	"pieces":      func(f float64) (string, float64) { return pcs, f },
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
func textToIngrds(s string) []Ingrd {
	lines := textToLines(s)
	ingrs := []Ingrd{}
	// Convert each line to an ingredient
	for _, line := range lines {
		var ingr Ingrd
		xs := strings.Split(line, " ")
		// Parse each element of the line to a float to find the amount
		for i, s := range xs {
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
				_, ok := unitsconv[xs[i+1]]
				if ok {
					unit, amount = unitsconv[xs[i+1]](amount)
					offset += len(xs[i+1])
				}
				ingr = Ingrd{
					Amount: amount,
					Unit:   unit,
					Item:   strings.Trim(line[strings.Index(line, s)+len(s)+offset+1:], " "), // assuming item is directly after amount in the text
					Notes:  strings.Trim(line[:strings.Index(line, s)], " "),                 // assuming an text before the float is additional notes
				}
				break
			}
		}
		if ingr.Amount == 0 && ingr.Unit == "" && ingr.Item == "" {
			// No amount found
			ingr.Item = line
		}
		ingrs = append(ingrs, ingr)
	}
	return ingrs
}
