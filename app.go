package main

import (
	"log"
	"os"
)

// Folders and file names used for config.
var (
	folderConfig   = "./config/"
	fnameRcps      = folderConfig + "recipes.json"
	fnameConvTable = folderConfig + "conversion.json"
	folderLog      = "./log/"
	fnameLog       = folderLog + "logfile.log"
)

/**checkFolder checks if folder f exists and if not creates the folder.*/
func checkFolder(f string) {
	if _, err := os.Stat(f); os.IsNotExist(err) {
		os.Mkdir(f, 4096)
	}
}

func init() {
	checkFolder(folderLog)
	checkFolder(folderConfig)
}

func main() {
	// Open/create logfile
	f, err := os.OpenFile(fnameLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Panic("Error setting log file:", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println("--------Start of program--------")

	// Load recipes
	err = readJSON(&rcps, fnameRcps)
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
