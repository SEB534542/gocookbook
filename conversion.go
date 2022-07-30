package main

import (
	"fmt"
	"math"
	"strings"
)

// Conv contains the item conversion from 1 gram to ml
var convTable = map[string]float64{}

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
}

// TODO: review and update conversions
var (
	tbspToMl = 14.7867648 // ml for 1 tablespoon.
	tspToMl  = 4.92892159 // ml for 1 teaspoon.
	cuptoMl  = 236.588237 // ml for 1 cup.
)

/*Uoms takes a pointer to an ingredient, determines the amount for two
pre-determined alternative Unit of Measurements combined into one string,
and updates this in the 'Alt' of the ingredient i.*/
func (i *Ingrd) uoms() {
	var xs []string
	switch i.Unit {
	case gram:
		m := round(gramToMl(i.Item, i.Amount))
		if m != 0.0 {
			c := round(m / cuptoMl)
			xs = append(xs, fmt.Sprintf("%v %v", c, cup), fmt.Sprintf("%v %v", m, ml))
		}
		// c toevoegen?
	case cup:
		m := round(i.Amount * cuptoMl)
		if m != 0.0 {
			xs = append(xs, fmt.Sprintf("%v %v", m, ml))
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

func gramToMl(item string, x float64) float64 {
	if f, ok := convTable[item]; ok {
		return x * f
	}
	return 0.0
}

func mlToGram(item string, x float64) float64 {
	if f, ok := convTable[item]; ok {
		return x / f
	}
	return 0.0
}

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

func round(f float64) float64 {
	return math.Round(f*1000) / 1000
}
