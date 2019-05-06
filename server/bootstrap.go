// reading and loading configuration
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func loadConfig(config *Config) {
	// Read and load configurtion json
	if file, fileErr := os.Open("./config.json"); fileErr != nil {
		log.Fatalf("Could not read ./config.json file : %s \n", fileErr.Error())
	} else {
		// Reading json to config variable
		configDecoder := json.NewDecoder(file)
		if decodingErr := configDecoder.Decode(config); decodingErr != nil {
			log.Fatalf("Failed to decode config : %s \n", decodingErr.Error())
		}
	}
}

func serveWeb(address string) {

	// start web on spacific address
	router := mux.NewRouter()
	log.Println("Listening on ", address)

	// slingpages resource routes
	router.PathPrefix("/sl-res/").
		Handler(http.StripPrefix("/sl-res/",
			http.FileServer(http.Dir("./admin/public"))))

	// admin routes
	adminRoutes(router)

	// Loading frontend routes ...
	// var routes []SlingRoute
	// if err := db.All(&routes); err != nil {
	// 	log.Fatalln("failed to load routes : ", err.Error())
	// } else {
	// 	for _, route := range routes {
	// 		mountRoute(route, router)
	// 	}
	// }

	router.HandleFunc(`/{rest:[a-zA-Z0-9=\-\/]*}`, func(w http.ResponseWriter, r *http.Request) {
		var slingRoute SlingRoute

		log.Println("Handling at rest routes")

		log.Println("PageURL for page ", r.URL.Path)
		if err := db.One("PageURL", r.URL.Path, &slingRoute); err == nil {
			log.Println("PageURL for page ", r.URL.Path)
			// var slingPage SlingPage
			// if err := db.Find("PageNumber", slingRoute.PageNumber, &slingPage); err == nil {
			renderPage(w, slingRoute)
			// } else {
			// log.Println("could not fetch page  for ", slingRoute.PageNumber)
			// }
		} else {
			log.Println("could not fetch route ", r.URL.Path)
		}

	})

	if httpErr := http.ListenAndServe(address, router); httpErr != nil {
		log.Fatalf("Failed to start web service : %s", httpErr.Error())
	} else {
		log.Println("Web service loaded ")
	}

}

func mountRoute(route SlingRoute, router *mux.Router) {
	router.HandleFunc(route.PageURL,
		func(w http.ResponseWriter, r *http.Request) {
			log.Print(r.URL.Path)
			renderPage(w, route)
		})
}
