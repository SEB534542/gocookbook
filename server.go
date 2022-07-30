package main

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	//	"io/ioutil"
	// "github.com/satori/go.uuid"
	// "golang.org/x/crypto/bcrypt"
)

var (
	tpl        *template.Template
	fm         = template.FuncMap{"fdateHM": hourMinute, "fsliceString": sliceToString, "fminutes": minutes, "fseconds": seconds}
	dbSessions = map[string]string{}
)

var (
	maxIngrs = 20 // Maximum amount of Ingredients that can be added on webpage.
	maxSteps = 20 // Maximum amount of Steps that can be added on webpage.
	convRows = 10 // Rows where additional conversion data can be added.
)

func init() {
	//Loading gohtml templates
	tpl = template.Must(template.New("").Funcs(fm).ParseGlob("./templates/*"))
}

func startServer(port int) {
	if port == 0 {
		port = 8081
		log.Printf("No port configured, using port %v", port)
	}
	// TODO: configure TSL
	cert := ""
	key := ""
	log.Printf("Launching website at localhost:%v", port)
	http.HandleFunc("/", handlerMain)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/recipe/", handlerRecipe)
	http.HandleFunc("/edit/", handlerEditRcp)
	http.HandleFunc("/add", handlerAddRcp)
	http.HandleFunc("/conv", handlerConversion)
	// http.HandleFunc("/log/", handlerLog)
	// http.HandleFunc("/login", handlerLogin)
	// http.HandleFunc("/logout", handlerLogout)
	// http.HandleFunc("/stop", handlerStop)
	err := http.ListenAndServeTLS(":"+fmt.Sprint(port), cert, key, nil)
	if err != nil {
		log.Printf("Unable to launch TLS, launching without TLS (%v)", err)
		log.Fatal(http.ListenAndServe(":"+fmt.Sprint(port), nil))
	}
}

func hourMinute(t time.Time) string {
	return t.Format("15:04")
}

func minutes(d time.Duration) string {
	return fmt.Sprint(d.Minutes())
}

func seconds(d time.Duration) string {
	return fmt.Sprint(d.Seconds())
}

func sliceToString(xs []string) string {
	return strings.Join(xs, ",")
}

func reverseXSS(xxs [][]string) [][]string {
	r := [][]string{}
	for i, _ := range xxs {
		r = append(r, xxs[len(xxs)-1-i])
	}
	return r
}

func reverseXS(xs []string) []string {
	r := []string{}
	for i, _ := range xs {
		r = append(r, xs[len(xs)-1-i])
	}
	return r
}

// StoTime receives a string of time (format hh:mm) and a day offset, and returns a type time with today's and the supplied hours and minutes + the offset in days
func stoTime(t string, days int) (time.Time, error) {
	timeNow := time.Now()
	timeHour, err := strconv.Atoi(t[:2])
	if err != nil {
		return time.Time{}, err
	}
	timeMinute, err := strconv.Atoi(t[3:])
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day()+days, int(timeHour), int(timeMinute), 0, 0, time.Local), nil
}

func readCSV(file string) [][]string {
	// Read the file
	f, err := os.Open(file)
	if err != nil {
		f, err := os.Create(file)
		if err != nil {
			log.Fatal("Unable to create csv", err)
		}
		f.Close()
		return [][]string{}
	}
	defer f.Close()
	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return lines
}

// Append CSV takes a filename and adds the new lines to the corresponding CSV file
func appendCSV(file string, newLines [][]string) {
	lines := readCSV(file)
	lines = append(lines, newLines...)
	// Write the file
	f, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}
	w := csv.NewWriter(f)
	if err = w.WriteAll(lines); err != nil {
		log.Fatal(err)
	}
}

// strToInt transforms string to an int and returns a positive int or zero
func strToInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if i < 0 {
		return 0, err
	}
	return i, err
}

// GetIP gets a requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
func getIP(req *http.Request) string {
	forwarded := req.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return req.RemoteAddr
}

// MaxIntSlice receives variadic parameter of integers and return the highest integer
func MaxIntSlice(xi ...int) int {
	var max int
	for i, v := range xi {
		if i == 0 || v > max {
			max = v
		}
	}
	return max
}

func stringToSlice(s string) []string {
	xs := strings.Split(s, ",")
	for i, v := range xs {
		xs[i] = strings.Trim(v, " ")
	}
	return xs
}

func handlerMain(w http.ResponseWriter, req *http.Request) {
	// TODO: add login functionality?
	data := struct {
		Recipes []Recipe
	}{
		rcps,
	}

	err := tpl.ExecuteTemplate(w, "index.gohtml", data)
	if err != nil {
		log.Fatalln(err)
	}
}

/* handlerRecipe determines the recipe ID, gathers the coresponding recipe and
if no. of persons is send along (through post method), the recipe is adjusted to
the new number of persons and it sends the response back.*/
func handlerRecipe(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.URL.Path[len("/recipe/"):])
	if err != nil {
		http.Redirect(w, req, "/", http.StatusBadRequest)
		return
	}
	rcp, err := findRecipe(rcps, id)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusNotFound)
		return
	}
	if req.Method == http.MethodPost {
		if persons, err := strconv.Atoi(req.PostFormValue("Persons")); err == nil {
			rcp = adjustRcp(rcp, persons)
		}
	}
	// Include/update alternate UOMs
	for i, _ := range rcp.Ingrs {
		rcp.Ingrs[i].uoms()
	}
	err = tpl.ExecuteTemplate(w, "recipe.gohtml", rcp)
	if err != nil {
		log.Fatalln(err)
	}
}

func handlerAddRcp(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		rcp := processRcp(req)
		rcp.Id = newRcpId(rcps)
		rcps = append(rcps, rcp)
		SaveToJSON(rcps, fnameRcps)
		http.Redirect(w, req, fmt.Sprintf("recipe/%v", rcp.Id), http.StatusSeeOther)
		return
	}
	data := struct {
		Recipe
		CountIngrs []int
		CountSteps []int
		Units      []string
	}{
		Recipe{},
		rangeList(0, maxIngrs),
		rangeList(0, maxSteps),
		units,
	}
	err := tpl.ExecuteTemplate(w, "add.gohtml", data)
	if err != nil {
		log.Fatalln(err)
	}
}

/* handlerEditRcp determines the recipe ID, gathers pointer to the coresponding
recipe and if no. of persons is send along (through post method), the recipe is
adjusted to the new number of persons and it sends the response back.*/
func handlerEditRcp(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.URL.Path[len("/edit/"):])
	if err != nil {
		http.Redirect(w, req, "/", http.StatusBadRequest)
		return
	}
	rcp, err := findRecipeP(rcps, id)
	if err != nil {
		http.Redirect(w, req, "/", http.StatusNotFound)
		return
	}
	if req.Method == http.MethodPost {
		rcpNew := processRcp(req)
		*rcp = rcpNew
		SaveToJSON(rcps, fnameRcps)
		http.Redirect(w, req, fmt.Sprintf("/recipe/%v", rcp.Id), http.StatusSeeOther)
	}
	data := struct {
		Recipe
		CountIngrs []int
		CountSteps []int
		Units      []string
	}{
		*rcp,
		rangeList(len(rcp.Ingrs), maxIngrs),
		rangeList(len(rcp.Steps), maxSteps),
		units,
	}
	err = tpl.ExecuteTemplate(w, "edit.gohtml", data)
	if err != nil {
		log.Fatalln(err)
	}
}

/* rangeList takes a min and max and return the numbers in between as a
slice of int*/
func rangeList(min, max int) []int {
	x := make([]int, max-min)
	for i := 0; i < (max - min); i++ {
		x[i] = min + i
	}
	return x
}

func processRcp(req *http.Request) Recipe {
	rcp := Recipe{}
	rcp.Name = req.PostFormValue("Name")
	rcp.Notes = req.PostFormValue("Notes")
	rcp.Persons, _ = strconv.Atoi(req.PostFormValue("Persons"))
	// Ingredients
	rcp.Ingrs = []Ingrd{}
	for i := 0; i < maxIngrs; i++ {
		ingr := Ingrd{}
		amount, _ := strconv.ParseFloat(req.PostFormValue(fmt.Sprintf("Amount%v", i)), 64)
		if amount == 0.0 {
			continue
		}
		ingr.Amount = amount
		ingr.Unit = req.PostFormValue(fmt.Sprintf("Unit%v", i))
		ingr.Item = strings.ToLower(req.PostFormValue(fmt.Sprintf("Item%v", i))) // All items are stored in lowercase.
		ingr.Notes = req.PostFormValue(fmt.Sprintf("Notes%v", i))
		rcp.Ingrs = append(rcp.Ingrs, ingr)
	}
	// Steps
	rcp.Steps = []string{}
	for i := 0; i < maxSteps; i++ {
		step := req.PostFormValue(fmt.Sprintf("Step%v", i))
		if step == "" {
			continue
		}
		rcp.Steps = append(rcp.Steps, step)
	}
	rcp.Source = req.PostFormValue("Source")
	return rcp
}

func handlerConversion(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		for k, _ := range convTable {
			if req.PostFormValue(fmt.Sprintf("%v-delete", k)) != "" {
				delete(convTable, k)
			} else {
				convTable[k], _ = strconv.ParseFloat(req.PostFormValue(k), 64)
			}
		}
		for i := 0; i < convRows; i++ {
			if k := strings.ToLower(req.PostFormValue(fmt.Sprint(i))); k != "" {
				convTable[k], _ = strconv.ParseFloat(req.PostFormValue(fmt.Sprintf("value-%v", i)), 64)
			}
		}
		SaveToJSON(convTable, fnameConvTable)
	}

	data := struct {
		ConvTable map[string]float64
		AddRows   []int
	}{
		convTable,
		rangeList(0, convRows),
	}
	err := tpl.ExecuteTemplate(w, "conversion.gohtml", data)
	if err != nil {
		log.Fatalln(err)
	}
}
