package main

import (
	"bytes"
	"encoding/csv"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

/* SaveToJson takes an interface and stores it into the filename*/
func SaveToJSON(i interface{}, fileName string) {
	bs, err := json.Marshal(i)
	if err != nil {
		log.Fatal(err)
	}
	//  Use below if you want JSON pretty printed
	var prettyJSON bytes.Buffer
	_ = json.Indent(&prettyJSON, bs, "", "    ")

	err = ioutil.WriteFile(fileName, prettyJSON.Bytes(), 0644)
	if err != nil {
		log.Fatal("Error saving JSON:", err)
	}
}

/*
	ReadJSON takes a pointer to an interface, reads from the given json file

location, stores in i and returns any error.
*/
func readJSON(i interface{}, fname string) error {
	if _, err := os.Stat(fname); os.IsNotExist(err) {
		return fmt.Errorf("File '%v' does not exist, creating new", fname)
		SaveToJSON(i, fname)
	} else {
		data, err := ioutil.ReadFile(fname)
		if err != nil {
			return fmt.Errorf("%s is corrupt. Please correct or delete the file (%v)", fname, err)
		}
		err = json.Unmarshal(data, &i)
		if err != nil {
			return fmt.Errorf("%s is corrupt. Please correct or delete the file (%v)", fname, err)
		}
	}
	return nil
}

/*
	ReadGob takes a pointer to an interface, reads from the given g fob ile

location, stores in i and returns any error.
*/
func ReadGob(i interface{}, fname string) error {
	// Initialize decoder
	var data bytes.Buffer
	dec := gob.NewDecoder(&data) // Will decode (read) and store into data

	// Read content from file
	content, err := ioutil.ReadFile(fname)
	if err != nil {
		return fmt.Errorf("Error reading file '%v': %v", fname, err)
	}
	y := bytes.NewBuffer(content)
	data = *y

	// Decode (receive) and print the values.
	err = dec.Decode(i)
	if err != nil {
		return fmt.Errorf("Error decoding into '%v': %v (%v)", fname, err, i)
	}
	return nil
}

// SaveGob encodes an interface and stores it as a Gob into a file named fname.
func SaveToGob(i interface{}, fname string) error {
	var data bytes.Buffer

	enc := gob.NewEncoder(&data) // Will write to data

	// Encode (send) some values.
	err := enc.Encode(i)
	if err != nil {
		return fmt.Errorf("Error encoding '%v': %v", fname, err)
	}

	// Store data
	err = ioutil.WriteFile(fname, data.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("Error storing '%v': %v", fname, err)
	}
	return nil
}

/*
	jsonStringPretty takes an interface and returns it as a string containing the

JSON structure for that interface pretty printed.
*/
func jsonStringPretty(i interface{}) (string, error) {
	bs, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	// Return string(bs) for json without pretty print.
	//  Use below and return prettyJSON.String if you want JSON pretty printed
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, bs, "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

/*
	jsonString takes an interface and returns it as a string containing the

JSON structure for that interface.
*/
func jsonString(i interface{}) (string, error) {
	bs, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

/*
	AppendCSV takes a filename and new lines and adds the new lines to the

corresponding CSV file.
*/
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

/*
	readCSV takes a filename to a CSV file and returns the CSV as a [][]string,

where the first slice represents each row and the second the comma separated
text on that line.
*/
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
