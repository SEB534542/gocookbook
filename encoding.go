package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

/* SaveToJson should take an interface and stores it into the filename*/
func SaveToJSON(i interface{}, fileName string) {
	bs, err := json.Marshal(i)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(fileName, bs, 0644)
	if err != nil {
		log.Fatal("Error", err)
	}
}

/* ReadJSON takes a pointer to an interface, reads from the given json file
location, stores in i and returns any error.*/
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

/* ReadGob takes a pointer to an interface, reads from the given g fob ile
location, stores in i and returns any error.*/
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

func jsonString(i interface{}) (string, error) {
	bs, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	// Use below and return prettyJSON if you want JSON pretty printed
	// var prettyJSON bytes.Buffer
	// if err := json.Indent(&prettyJSON, bs, "", "    "); err != nil {
	// 	return "", err
	// }
	return string(bs), nil
}
