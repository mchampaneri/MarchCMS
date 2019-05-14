package main

import (
	"log"
	"os"

	"github.com/satori/go.uuid"

	"github.com/asdine/storm"
)

// Config holds global configurations
// of cms
type Config struct {
	ID       string `json:"id"`
	Address  string `json:"Address"`
	Name     string `json:"Name"`
	Database string `json:"Database"`
	Theme    string `json:"Theme"`
	Status   string `json:"Status"`
}

// Global db variables
var db *storm.DB
var dbErr error

// CMS wide config
var config, jsonConfig Config
var root, _ = os.Getwd()

func main() {

	// Loading configurations
	loadConfig(&jsonConfig)
	// Preapring db
	db, dbErr = storm.Open("my.db")
	// preparing fake data
	// feedFakeData()
	if err := db.One("Status", "Active", &config); err != nil {
		idForConfig, _ := uuid.NewV4()
		jsonConfig.ID = idForConfig.String()
		log.Println("Could not get local config", err.Error())
		if err := db.Save(&jsonConfig); err != nil {
			log.Fatalln("Could not save local config", err.Error())
		} else {
			if err := db.One("Status", "Active", &config); err != nil {
				log.Fatalln("Could not save local config", err.Error())
			}
		}
	}
	// loading webservice
	serveWeb(config.Address) // loading web service

	defer db.Close()
}
