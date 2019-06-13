package main

import (
	"log"
	"net/rpc"
	"os"
	"path/filepath"

	"github.com/asdine/storm"
	uuid "github.com/satori/go.uuid"
)

// Global db variables
var db *storm.DB
var dbErr error

// CMS wide config
var config, jsonConfig Config
var root, _ = os.Getwd()

// Extension Handles
var extensions = make(map[string]*rpc.Client)

// // Folder Paths
var themesFolder = filepath.Join(root, "themes")
var adminFolder = filepath.Join(root, "admin")
// var extensionFolder = filepath.Join(root, "extensions")

func main() {

	// Loading configurations
	loadConfig(&jsonConfig)

	// loadExtensions()
	// Preapring db
	db, dbErr = storm.Open("my.db")

	// Find and load active config from db
	// or prepare new one from json
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

	// loadExtensions()
	// loading webservice
	serveWeb(config.Address) // loading web service
	defer db.Close()
}
