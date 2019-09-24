package main

import (
	"flag"
	"log"
	"net/rpc"
	"os"
	"path/filepath"

	"github.com/asdine/storm"
	"github.com/gorilla/sessions"
	uuid "github.com/satori/go.uuid"
)

// Global db variables
var db *storm.DB
var dbErr error

// CMS wide config
var config, jsonConfig Config
var themeConfig, themeJsonConfig ThemeConfig
var root, _ = os.Getwd()

// Extension Handles
var extensions = make(map[string]*rpc.Client)
var UserSession = sessions.NewCookieStore([]byte("xf7KylXJ7CFSH4mLZG2Wyl86HAB9Rqvn"))

// // Folder Paths
var themesFolder = filepath.Join(root, "themes")
var adminFolder = filepath.Join(root, "admin")
var assetFolder = filepath.Join(root, "assets")

// var extensionFolder = filepath.Join(root, "extensions")

func main() {

	systemCheck()
	// Loading configurations
	loadSiteConfig(&jsonConfig)

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

	// Loading config of active theme
	loadThemeConfig(config.Theme, &themeJsonConfig)
	// Load theme settings
	if err := db.One("Theme", config.Theme, &themeConfig); err != nil {
		idForConfig, _ := uuid.NewV4()
		themeJsonConfig.ID = idForConfig.String()
		log.Println("Could not get local theme config", err.Error())
		if err := db.Save(&themeJsonConfig); err != nil {
			log.Fatalln("Could not save local theme config", err.Error())
		} else {
			if err := db.One("Theme", config.Theme, &themeConfig); err != nil {
				log.Fatalln("Could not get saved local theme config ", config.Theme, err.Error())
			}
		}
	}

	// Find and load active config from db
	// or prepare new one from json
	AdamUser := MarchUser{
		Name:      "Adam MarchCMS",
		Email:     "adam@marchcms.org",
		SmallDesc: "Adam is root user of MarchCMS",
		Role:      adminUser,
		Password:  "adam@1234",
		Status:    activeAccount,
	}

	if err := db.One("Email", "adam@marchcms.org", &AdamUser); err != nil {
		log.Println("Could not get adam user", err.Error())
		if user, err := AdamUser.RegisterUser(); err != nil {
			log.Println("Failed to generate Adam User : ", err.Error())
		} else {
			log.Println("Adam user is :", user)
		}
	}

	// Reading commmadn line inputs
	// It can change runtime value of config
	wordPtr := flag.String("live", "yes", "a string")

	flag.Parse()
	if *wordPtr == "Yes" {
		goLive()
	} else {
		devMode()
	}
	// loadExtensions()
	// loading webservice
	serveWeb(config.Address) // loading web service
	defer db.Close()
}
