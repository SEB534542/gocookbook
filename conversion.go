package main

import (
	"fmt"
	"math"
	"strings"
)

// Conv contains the item conversion from 100 gram to a cup
var convTable = map[string]float64{}

const (
	gram = "gr"
	cup  = "cup"
	ml   = "ml"
	tbsp = "el"
	tsp  = "tl"
	pcs  = "stuks"
)

var units = []string{
	gram, cup, ml, tbsp, tsp, pcs,
}

var (
	cupToMilliliter = 0.2841306 // fixed cup to milliliter ratio
	cupToTbsp       = 16.0      // fixed cup to tablespoon ration
	cupToTsp        = 48.0      // fixed cup to teaspoon ratio
)

/*Uoms takes a pointer to an ingredient, determines the amount for two
pre-determined alternative Unit of Measurements combined into one string,
and updates this in the 'Alt' of the ingredient i.*/
func (i *Ingrd) uoms() {
	var xs []string
	switch i.Unit {
	case gram:
		c := round(gramToCup(i.Item, i.Amount))
		if c != 0.0 {
			xs = append(xs, fmt.Sprintf("%v %v", c, cup))
		}
		m := round(c * cupToMilliliter)
		if m != 0.0 {
			xs = append(xs, fmt.Sprintf("%v %v", m, ml))
		}
	case cup:
		g := round(cupToGram(i.Item, i.Amount))
		if g != 0.0 {
			xs = append(xs, fmt.Sprintf("%v %v", g, gram))
		}
		m := round(i.Amount * cupToMilliliter)
		if m != 0.0 {
			xs = append(xs, fmt.Sprintf("%v %v", m, ml))
		}
	case ml:
		c := round(1 / cupToMilliliter * i.Amount)
		if c != 0.0 {
			xs = append(xs, fmt.Sprintf("%v %v", c, cup))
		}
		g := round(cupToGram(i.Item, c))
		if g != 0.0 {
			xs = append(xs, fmt.Sprintf("%v %v", g, gram))
		}
	case tbsp:
		c := round(1 / cupToTbsp * i.Amount)
		if c != 0.0 {
			xs = append(xs, fmt.Sprintf("%v %v", c, cup))
		}
		g := cupToGram(i.Item, c)
		if g != 0.0 {
			xs = append(xs, fmt.Sprintf("%v %v", g, gram))
		}
	case tsp:
		c := round(1 / cupToTsp * i.Amount)
		if c != 0.0 {
			xs = append(xs, fmt.Sprintf("%v %v", c, cup))
		}
		g := cupToGram(i.Item, c)
		if g != 0.0 {
			xs = append(xs, fmt.Sprintf("%v %v", g, gram))
		}
	}
	i.AltUnits = strings.Join(xs, " / ")
}

func gramToCup(item string, x float64) float64 {
	if f, ok := convTable[item]; ok {
		return x / 100 * f
	}
	return 0.0
}

func cupToGram(item string, x float64) float64 {
	if f, ok := convTable[item]; ok {
		return x * 1 / f * 100
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
