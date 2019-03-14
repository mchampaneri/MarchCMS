package main

import (
	"github.com/asdine/storm"
)

// Config holds global configurations
// of cms
type Config struct {
	Address  string `json:"Address"`
	Name     string `json:"Name"`
	Database string `json:"Database"`
}

// Global db variables
var db *storm.DB
var dbErr error

// CMS wide config
var config Config

func main() {

	// Loading configurations
	loadConfig(&config)
	// Preapring db
	db, dbErr = storm.Open("my.db")
	// preparing fake data
	// feedFakeData()

	// loading webservice
	serveWeb(config.Address) // loading web service

	defer db.Close()
}
