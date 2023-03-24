package main

import (
	"fmt"
	"math"
	"strings"
)

// convTable contains the item conversion from 1 gram to ml.
var convTable = map[string]float64{}

/*
	different types of volumes and masses used for conversion. Note don't change

the actual string without changing the existing data and adding it to the var
units.
*/
const (
	gram = "g"
	cup  = "cup"
	ml   = "ml"
	tbsp = "el"
	tsp  = "tl"
	pcs  = "stuks"
)

var units = []string{
	gram, cup, ml, tbsp, tsp, pcs,
} // all considered volumes and masses that are used in the cookbook.

var (
	tbspToMl = 14.7867648 // ml for 1 tablespoon.
	tspToMl  = 4.92892159 // ml for 1 teaspoon.
	cuptoMl  = 236.588237 // ml for 1 cup.
)

/*
Uoms takes a pointer to an ingredient, determines the amount for two
pre-determined alternative Unit of Measurements combined into one string,
and updates this in the 'Alt' of the ingredient i.
*/
func (i *Ingrd) uoms() {
	var xs []string
	switch i.Unit {
	case gram:
		m := round(gramToMl(i.Item, i.Amount))
		if m != 0.0 {
			c := round(m / cuptoMl)
			xs = append(xs, fmt.Sprintf("%v %v", m, ml), fmt.Sprintf("%v %v", c, cup))
		}
	case cup:
		m := round(i.Amount * cuptoMl)
		if m != 0.0 {
			xs = append(xs, fmt.Sprintf("%v %v", m, ml))
			g := round(mlToGram(i.Item, m))
			if g != 0.0 {
				xs = append(xs, fmt.Sprintf("%v %v", g, gram))
			}
		}
	case ml:
		c := round(i.Amount / cuptoMl)
		xs = append(xs, fmt.Sprintf("%v %v", c, cup))
		g := round(mlToGram(i.Item, i.Amount))
		if g != 0.0 {
			xs = append(xs, fmt.Sprintf("%v %v", g, gram))
		}
	case tbsp:
		m := round(i.Amount * tbspToMl)
		c := round(1 / cuptoMl * m)
		xs = append(xs, fmt.Sprintf("%v %v", m, ml), fmt.Sprintf("%v %v", c, cup))
		g := mlToGram(i.Item, m)
		if g != 0.0 {
			xs = append(xs, fmt.Sprintf("%v %v", g, gram))
		}
	case tsp:
		m := round(i.Amount * tspToMl)
		c := round(1 / cuptoMl * m)
		xs = append(xs, fmt.Sprintf("%v %v", m, ml), fmt.Sprintf("%v %v", c, cup))
		g := mlToGram(i.Item, m)
		if g != 0.0 {
			xs = append(xs, fmt.Sprintf("%v %v", g, gram))
		}
	}
	i.AltUnits = strings.Join(xs, " / ")
}

/*
	gramToMl takes an item and number of grams, looks up the item in the

conversion table and returns the number of milliliters for x grams of the item.
*/
func gramToMl(item string, x float64) float64 {
	if f, ok := convTable[item]; ok {
		return x * f
	}
	return 0.0
}

/*
	mlToGram takes an item and number of milliliters, looks up the item in the

conversion table and returns the number of grams for x milliliters of the item.
*/
func mlToGram(item string, x float64) float64 {
	if f, ok := convTable[item]; ok {
		return x / f
	}
	return 0.0
}

/*
	toTitle takes a string, capitalizes the first value and sets the rest to lower

case.
*/
func toTitle(s string) string {
	if len(s) == 0 {
		return ""
	}
	newS := fmt.Sprint(strings.ToUpper(string(s[0])))
	if len(s) > 1 {
		newS += strings.ToLower(s[1:])
	}
	return newS
}

/*
	round takes a float and rounds it to three decimals. E.g. round(0.5555) returns

0.555.
*/
func round(f float64) float64 {
	return math.Round(f*10) / 10
}
