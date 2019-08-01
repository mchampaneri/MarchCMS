package main

import "log"

// Makes Live status of current site config
// to "yes"
func goLive() {
	config.Live = "Yes"
	log.Println("Making site live")
}

// Makes Live status of current site config
// to "no"
func devMode() {
	config.Live = "No"
	log.Println("Making site under maintainence")
}
