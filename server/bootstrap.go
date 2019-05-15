// reading and loading configuration
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
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

func loadExtensions() {
	client, err := rpc.Dial("tcp", "127.0.0.1:10000")
	defer client.Close()
	if err != nil {
		log.Fatal(err.Error())
	}

	response := Response{}
	callErr := client.Call("HelloRequest.Hello",
		Request{Name: "Manish"},
		&response)

	if callErr != nil {
		log.Fatalln("failed to make call :", callErr.Error())
	}

	fmt.Println("responsoe from extnsion is below")
	fmt.Println(response.Message)
}

func serveWeb(address string) {

	// start web on spacific address
	router := mux.NewRouter()
	log.Println("Listening on ", address)

	// slingpages resource routes
	router.PathPrefix("/sl-res/").
		Handler(http.StripPrefix("/sl-res/",
			http.FileServer(http.Dir("./admin/public"))))

	// theme resource routes
	router.PathPrefix("/th-res/").
		Handler(http.StripPrefix("/th-res/",
			http.FileServer(http.Dir("./themes/shortshot/public"))))

	// admin routes
	adminRoutes(router)

	// Handling posts
	router.HandleFunc(`/post/{rest:[a-zA-Z0-9=\-\/]*}`, func(w http.ResponseWriter, r *http.Request) {
		var slingpost SlingPost
		log.Println("Handling at rest routes")

		log.Println("PageURL for page ", r.URL.Path)
		if err := db.One("PageURL", r.URL.Path, &slingpost); err == nil {
			log.Println("PageURL for page ", r.URL.Path)
			renderPost(w, slingpost)
		} else {
			log.Println("could not fetch route ", r.URL.Path)
		}
	})

	// Handling pages
	router.HandleFunc(`/{rest:[a-zA-Z0-9=\-\/]*}`, func(w http.ResponseWriter, r *http.Request) {
		var slingPage SlingPage
		log.Println("Handling at rest routes")

		log.Println("PageURL for page ", r.URL.Path)
		if err := db.One("PageURL", r.URL.Path, &slingPage); err == nil {
			log.Println("PageURL for page ", r.URL.Path)
			renderPage(w, slingPage)
		} else {
			renderPage(w, SlingPage{PageTemplate: "404.html"})
		}
	})

	if httpErr := http.ListenAndServe(address, router); httpErr != nil {
		log.Fatalf("Failed to start web service : %s", httpErr.Error())
	} else {
		log.Println("Web service loaded ")
	}

}
