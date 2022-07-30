package main

import (
	"log"
	"os"
)

var (
	folderConfig   = "./config/"
	fnameRcps      = folderConfig + "recipes.json"
	fnameConvTable = folderConfig + "conversion.json"
)

func main() {
	// Check if log folder exists, else create
	if _, err := os.Stat(folderConfig); os.IsNotExist(err) {
		os.Mkdir(folderConfig, 4096)
	}
	// Load recipes
	err := readJSON(&rcps, fnameRcps)
	if err != nil {
		log.Println(err)
	}
	// Load conversion table
	err = readJSON(&convTable, fnameConvTable)
	if err != nil {
		log.Println(err)
	}
	startServer(8081)
}
