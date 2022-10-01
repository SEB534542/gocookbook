package main

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Username string
	Password []byte
}

type visit struct {
	Ip   string    // IP address origin
	Port string    // Port
	Time time.Time // Datetime of visit
	Site string    // Actual site visited
	Un   string    // If logged in, username of visitor
}

// Folders and file names used
var (
	fnameUsers  = folderConfig + "users.json"
	fnameVisits = folderLog + "visits.json"
)

var (
	tpl *template.Template
	fm  = template.FuncMap{
		"fdateHM":           hourMinute,
		"fsliceString":      sliceToString,
		"fsliceStringSpace": sliceToStringSpace,
		"fminutes":          minutes,
		"fseconds":          seconds,
		"fdate":             dateTime,
		"fplusOne":          plusOne,
	} // Map with all functions that can be used within html.
	dbUsers    = map[string]user{}   // username, user
	dbSessions = map[string]string{} // session ID, username
	dbVisits   = []visit{}           // Visits to this website.
)

var (
	maxIngrs = 30 // Maximum amount of Ingredients that can be added on webpage.
	maxSteps = 20 // Maximum amount of Steps that can be added on webpage.
	convRows = 10 // Rows where additional conversion data can be added.
)

const cookieSession = "session"

func init() {
	//Loading gohtml templates
	tpl = template.Must(template.New("").Funcs(fm).ParseGlob("./templates/*"))
}

/* startServer takes a port and launches a server. It tries to create a HTTPS
server, but if that fails, it creates a HTTP server.*/
func startServer(port int) {
	if port == 0 {
		port = 8081
		log.Printf("No port configured, using port %v", port)
	}
	// load users
	err := readJSON(&dbUsers, fnameUsers)
	if err != nil {
		log.Printf("Unable to load users from '%v': %v", fnameUsers, err)
	}
	// load visits
	err = readJSON(&dbVisits, fnameVisits)
	if err != nil {
		log.Printf("Unable to load previous visits from '%v': %v", fnameVisits, err)
	}
	// TODO: configure TSL
	cert := "" // address of cert file
	key := ""  // address of key file
	log.Printf("Launching website at localhost:%v", port)
	http.HandleFunc("/", handlerMain)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/recipe/", handlerRecipe)
	http.HandleFunc("/edit/", handlerEditRcp)
	http.HandleFunc("/add", handlerAddRcp)
	http.HandleFunc("/conv", handlerConversion)
	http.HandleFunc("/export/recipes", handlerExportRcps)
	http.HandleFunc("/export/table", handlerExportTable)
	http.HandleFunc("/log/", handlerLog)
	http.HandleFunc("/login", handlerLogin)
	http.HandleFunc("/logout", handlerLogout)
	http.HandleFunc("/visits", handlerVisits)
	err = http.ListenAndServeTLS(":"+fmt.Sprint(port), cert, key, nil)
	if err != nil {
		log.Printf("Unable to launch TLS, launching without TLS (%v)", err)
		log.Fatal(http.ListenAndServe(":"+fmt.Sprint(port), nil))
	}
}

// hourMinute takes a time.Time and returns it as a string.
func hourMinute(t time.Time) string {
	return t.Format("15:04")
}

// dateTime takes a time.Time and returns it as a string.
func dateTime(t time.Time) string {
	return t.Format("02-01 15:04")
}

// minutes takes a duration and returns the minutes as a string.
func minutes(d time.Duration) string {
	return fmt.Sprint(d.Minutes())
}

// seconds takes a duration and returns the seconds as a string.
func seconds(d time.Duration) string {
	return fmt.Sprint(d.Seconds())
}

// sliceToString takes a slice of string and returns it is a string, with each value comma separated.
func sliceToString(xs []string) string {
	return strings.Join(xs, ",")
}

// sliceToStringSpace takes a slice of string and returns it is a string, with each value comma+space separated.
func sliceToStringSpace(xs []string) string {
	return strings.Join(xs, ", ")
}

// reverseXSS takes a slice of a slice of string and returns it in reversed order.
func reverseXSS(xxs [][]string) [][]string {
	r := [][]string{}
	for i, _ := range xxs {
		r = append(r, xxs[len(xxs)-1-i])
	}
	return r
}

// reverseXS takes a slice of string and reutnrs it in reversed order.
func reverseXS(xs []string) []string {
	r := []string{}
	for i, _ := range xs {
		r = append(r, xs[len(xs)-1-i])
	}
	return r
}

/* StoTime receives a string of time (format hh:mm) and a day offset, and
returns a type time with today's and the supplied hours and minutes + the offset
in days.*/
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

// strToInt transforms string to an int and returns a positive int or zero.
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

/*CheckIp takes a map of IP addresses and an IP address, checks if the
address is already present in the map and stores this in the log. If the address
is local (i.e. starts with 192), it omits the address from the log.*/
func checkIp(ips map[string]bool, ip string) {
	if _, ok := ips[ip]; !ok {
		ips[ip] = true
		if ip[:3] != "192" {
			log.Printf("New ip visited: '%v'", ip)
		}
	}
}

/* GetIP takes a request's IP address by reading off the forwarded-for
header (for proxies) and returns the to use the remote address.*/
func getIP(req *http.Request) string {
	forwarded := req.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return req.RemoteAddr
}

/* handlerVisits prints all stored visits on a HTML page.*/
func handlerVisits(w http.ResponseWriter, req *http.Request) {
	addVisit(req)
	if !alreadyLoggedIn(req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	err := tpl.ExecuteTemplate(w, "visits.gohtml", dbVisits)
	if err != nil {
		log.Fatalln(err)
	}
}

/*handlerLog displays the complete log.*/
func handlerLog(w http.ResponseWriter, req *http.Request) {
	addVisit(req)
	if !alreadyLoggedIn(req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	f, err := ioutil.ReadFile(fnameLog)
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	lines := strings.Split(string(f), "\n")

	var output string
	for _, v := range lines {
		output += fmt.Sprintln(v)
	}
	fmt.Fprintf(w, output)
}

/* MaxIntSlice receives variadic parameter of integers and return the highest
integer.*/
func maxIntSlice(xi ...int) int {
	var max int
	for i, v := range xi {
		if i == 0 || v > max {
			max = v
		}
	}
	return max
}

/* stringToSlice takes a string and returns a slice of string, for each comma.*/
func stringToSlice(s string) []string {
	xs := strings.Split(s, ",")
	for i, v := range xs {
		xs[i] = strings.Trim(v, " ")
	}
	return xs
}

func handlerMain(w http.ResponseWriter, req *http.Request) {
	addVisit(req)
	data := struct {
		Recipes []Recipe
		Tags    []string
		Known   bool
	}{
		rcps,
		tags(rcps),
		alreadyLoggedIn(req),
	}

	err := tpl.ExecuteTemplate(w, "index.gohtml", data)
	if err != nil {
		log.Fatalln(err)
	}
}

/* handlerExportRcps prints all recipes in JSON on the webpage.*/
func handlerExportRcps(w http.ResponseWriter, req *http.Request) {
	addVisit(req)
	if !alreadyLoggedIn(req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	output, err := jsonStringPretty(rcps)
	if err != nil {
		msg := "Error saving:" + fmt.Sprint(err)
		http.Error(w, msg, http.StatusExpectationFailed)
	}
	fmt.Fprintf(w, output)
}

/* handlerExportTable prints the conversion table in JSON on the webpage.*/
func handlerExportTable(w http.ResponseWriter, req *http.Request) {
	addVisit(req)
	if !alreadyLoggedIn(req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	output, err := jsonStringPretty(convTable)
	if err != nil {
		msg := "Error saving:" + fmt.Sprint(err)
		http.Error(w, msg, http.StatusExpectationFailed)
	}
	fmt.Fprintf(w, output)
}

/* handlerRecipe determines the recipe ID, gathers the coresponding recipe and
if no. of persons is send along (through post method), the recipe is adjusted to
the new number of persons and it sends the response back.*/
func handlerRecipe(w http.ResponseWriter, req *http.Request) {
	addVisit(req)
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
	data := struct {
		Recipe Recipe
		Known  bool
	}{
		rcp,
		alreadyLoggedIn(req),
	}
	err = tpl.ExecuteTemplate(w, "recipe.gohtml", data)
	if err != nil {
		log.Fatalln(err)
	}
}

/* handlerAddRcp generates the html page to enter a new recipe and processes and
stores the new recipe.*/
func handlerAddRcp(w http.ResponseWriter, req *http.Request) {
	addVisit(req)
	if !alreadyLoggedIn(req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	if req.Method == http.MethodPost {
		rcp := processRcp(req)
		rcp.Id = newRcpId(rcps)
		rcps = append(rcps, rcp)
		sort.Slice(rcps, func(i, j int) bool { return rcps[i].Name < rcps[j].Name })
		SaveToJSON(rcps, fnameRcps)
		log.Printf("New recipe added (id %v", rcp.Id)
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

/* handlerEditRcp lookus up the recipe ID from the path, generates the recipe
on the html page and processes any updates.*/
func handlerEditRcp(w http.ResponseWriter, req *http.Request) {
	addVisit(req)
	if !alreadyLoggedIn(req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
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
		sort.Slice(rcps, func(i, j int) bool { return rcps[i].Name < rcps[j].Name })
		SaveToJSON(rcps, fnameRcps)
		log.Printf("Recipe %v updated", rcpNew.Id)
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

/* processRcp takes a *http.requested and extracts the form POST data into a
recipe, which is returned.*/
func processRcp(req *http.Request) Recipe {
	rcp := Recipe{}
	if id := req.PostFormValue("Id"); id != "" {
		rcp.Id, _ = strconv.Atoi(id)
	}
	rcp.Name = req.PostFormValue("Name")
	rcp.Notes = req.PostFormValue("Notes")
	rcp.Dur, _ = time.ParseDuration(fmt.Sprintf("%vm", req.PostFormValue("Dur")))
	rcp.Persons, _ = strconv.Atoi(req.PostFormValue("Persons"))

	t := stringToSlice(req.PostFormValue("Tags"))
	rcp.Tags = []string{}
	for _, v := range t {
		v = strings.Trim(v, " ")
		if v != "" {
			rcp.Tags = append(rcp.Tags, toTitle(v))
		}
	}
	sort.Strings(rcp.Tags)
	// Ingredients
	// Gather all ingredients
	ids := []float64{}
	ingrs := map[float64]Ingrd{}
	for i := 0; i < maxIngrs; i++ {
		ingr := Ingrd{}
		id, _ := strconv.ParseFloat(req.PostFormValue(fmt.Sprintf("Id%v", i)), 64)
		amount, _ := strconv.ParseFloat(req.PostFormValue(fmt.Sprintf("Amount%v", i)), 64)
		if amount == 0.0 {
			continue
		}
		ingr.Amount = amount
		ingr.Unit = req.PostFormValue(fmt.Sprintf("Unit%v", i))
		ingr.Item = strings.ToLower(req.PostFormValue(fmt.Sprintf("Item%v", i))) // All items are stored in lowercase.
		ingr.Notes = req.PostFormValue(fmt.Sprintf("Notes%v", i))
		ingrs[id] = ingr
		ids = append(ids, id)
	}
	// Sort and store ingredients into recipe
	rcp.Ingrs = make([]Ingrd, len(ingrs))
	sort.Float64s(ids)
	for i, id := range ids {
		rcp.Ingrs[i] = ingrs[id]
	}
	// Steps
	// Gather all steps
	stepIds := []float64{}
	steps := map[float64]string{}
	for i := 0; i < maxSteps; i++ {
		id, _ := strconv.ParseFloat(req.PostFormValue(fmt.Sprintf("StepId%v", i)), 64)
		step := req.PostFormValue(fmt.Sprintf("Step%v", i))
		if step == "" {
			continue
		}
		steps[id] = step
		stepIds = append(stepIds, id)
	}
	// Sort and store steps into recipe
	rcp.Steps = make([]string, len(stepIds))
	sort.Float64s(stepIds)
	for i, id := range stepIds {
		rcp.Steps[i] = steps[id]
	}
	rcp.Source = req.PostFormValue("Source")
	rcp.SourceLink = req.PostFormValue("SourceLink")
	switch {
	case rcp.SourceLink == "" && isHyperlink(rcp.Source):
		rcp.SourceLink = rcp.Source
	case rcp.Source == "" && isHyperlink(rcp.SourceLink):
		rcp.Source = rcp.SourceLink
	case !isHyperlink(rcp.SourceLink):
		rcp.SourceLink = ""
	}
	return rcp
}

/* handlerConversion generates the html page to show and update the
conversion table.*/
func handlerConversion(w http.ResponseWriter, req *http.Request) {
	addVisit(req)
	if !alreadyLoggedIn(req) {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
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

func handlerLogin(w http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	ip := getIP(req)
	// process form submission
	if req.Method == http.MethodPost {
		un := req.FormValue("Username")
		p := req.FormValue("Password")
		// Lookup username
		u, ok := dbUsers[un]
		if !ok {
			log.Printf("%v entered incorrect username %v..", un, ip)
			http.Error(w, "Username and/or password do not match", http.StatusForbidden)
			return
		}
		// Does the entered password match the stored password?
		err := bcrypt.CompareHashAndPassword(u.Password, []byte(p))
		if err != nil {
			log.Printf("%v entered incorrect password...", ip)
			http.Error(w, "Username and/or password do not match", http.StatusForbidden)
			return
		}
		// create session
		log.Printf("User (%v) logged in...", ip)
		sID := uuid.NewV4()
		c := &http.Cookie{
			Name:   "session",
			Value:  sID.String(),
			MaxAge: 0,
		}
		http.SetCookie(w, c)
		dbSessions[c.Value] = un
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	err := tpl.ExecuteTemplate(w, "login.gohtml", nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func handlerLogout(w http.ResponseWriter, req *http.Request) {
	addVisit(req)
	if !alreadyLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	c, _ := req.Cookie(cookieSession)
	// delete the session
	delete(dbSessions, c.Value)
	// remove the cookie
	c = &http.Cookie{
		Name:   cookieSession,
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func alreadyLoggedIn(req *http.Request) bool {
	c, err := req.Cookie(cookieSession)
	if err != nil {
		// Error retrieving cookie
		return false
	}
	un := dbSessions[c.Value]
	if _, ok := dbUsers[un]; !ok {
		// Unknown cookie and/or user
		return false
	}
	return true
}

/*username takes a http request, checks the session cookie to identify
and return the user (if).*/
func username(req *http.Request) string {
	c, err := req.Cookie(cookieSession)
	if err != nil {
		// Error retrieving cookie
		return ""
	}
	un := dbSessions[c.Value]
	if _, ok := dbUsers[un]; !ok {
		// Unknown cookie and/or user
		return ""
	}
	return un
}

/* isHyperlink takes a string, checks if a hyperlink exists in that string.*/
func isHyperlink(s string) bool {
	if s == "" {
		// string is empty, so not a hyperlink
		return false
	}
	xs := strings.Split(s, " ")
	if len(xs) > 1 {
		// string has spaces so cannot be a hyperlink
		return false
	}
	tags := []string{"http://", "https://", "www."}
	for _, t := range tags {
		if startsWith(s, t) {
			return true
		}
	}
	return false
}

/* startsWith takes a string and a substring and returns true if the string
starts with the substring. It does not require cases to match
(i.e. lower v/s upper case).*/
func startsWith(s, substr string) bool {
	if len(s) < len(substr) {
		return false
	}
	if strings.ToLower(s[:len(substr)]) == strings.ToLower(substr) {
		return true
	}
	return false
}

/* removeFromSlice takes a slice of string and removes the i-th element and
returns the array.*/
func removeFromSlice(xs []string, i int) []string {
	v := make([]string, len(xs)-1)
	v = append(xs[:i], xs[i+1:]...) // remove i-th element
	return v
}

func addVisit(req *http.Request) {
	ipp := getIP(req)
	site := req.URL.Path
	un := username(req)
	addr := strings.Split(ipp, ":") // ipp contains ip:port, ie 192.168.1.1:7000 and converts this into a slice of string.
	var ip, port string
	ip = addr[0]
	if len(addr) > 0 {
		port = addr[1]
	}
	v := visit{
		Ip:   ip,
		Port: port,
		Time: time.Now(),
		Site: site,
		Un:   un,
	}
	dbVisits = append(dbVisits, v)
	SaveToJSON(dbVisits, fnameVisits)
}

/* Increment takes an integer and returns the integer +1*/
func plusOne(i int) int {
	return i + 1
}

/*Tags receives a slice of Recipe and returns the unique Tags as a slice
of string*/
func tags(rcps []Recipe) []string {
	list := map[string]bool{}
	for _, rcp := range rcps {
		for _, tag := range rcp.Tags {
			list[tag] = true
		}
	}
	output := []string{}
	for k, _ := range list {
		output = append(output, k)
	}
	sort.Strings(output)
	return output
}

/*remove takes a slice of string and a position i and removes i from the
slice, by substituting it with the value at the end of the slice, as this is
less costly then moving all elements.*/
func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
