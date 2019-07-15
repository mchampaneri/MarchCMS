// reading and loading configuration
package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/handlers"
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

// func loadExtensions() {
// 	if fileInfo, err := ioutil.ReadDir(filepath.Join(root, "extensions")); err == nil {
// 		for _, file := range fileInfo {
// 			configFile := filepath.Join(extensionFolder, file.Name(), "config.json")
// 			readFile, err := os.Open(configFile)
// 			if err != nil {
// 				log.Println(err.Error())
// 			}
// 			var extensionConfig RpcExtension
// 			configDecoder := json.NewDecoder(readFile)
// 			configDecoder.Decode(&extensionConfig)

// 			if extensionConfig.Status == "active" {
// 				// if freePort, err := getAvailablePort(); err == nil {
// 				// Registring extension
// 				if client, err := rpc.Dial("tcp", extensionConfig.Address); err == nil {
// 					log.Println(extensionConfig.Name, " : loaded")
// 					extensions[extensionConfig.Name] = client
// 				} else {
// 					log.Println("Failed to register extension : ", extensionConfig.Name, err.Error())
// 				}
// 				// }
// 			}
// 		}
// 	} else {
// 		log.Fatalln("failed to load extension ", err.Error())
// 	}
// }

func serveWeb(address string) {

	// start web on spacific address
	router := mux.NewRouter()
	log.Println("Listening on ", address)

	// assets routes
	router.HandleFunc("/asset/uploaded/{ofType}/{asset}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ofType := vars["ofType"]
		requestedAssetName := vars["asset"]
		if file, err := os.Open(filepath.Join(assetFolder, ofType, requestedAssetName)); err == nil {
			io.Copy(w, file)
		} else {
			renderJSON(w, map[string]string{"error": "Could not find asset"})
		}
	})

	// marchPages resource routes
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
		var marchPost MarchPost
		log.Println("Handling at rest routes")

		log.Println("PageURL for page ", r.URL.Path)
		if err := db.One("PageURL", r.URL.Path, &marchPost); err == nil {
			log.Println("PageURL for page ", r.URL.Path)
			renderPost(w, r, marchPost)
		} else {
			log.Println("could not fetch route ", r.URL.Path)
		}
	})

	// Handling pages
	router.HandleFunc(`/{rest:[a-zA-Z0-9=\-\/]*}`, func(w http.ResponseWriter, r *http.Request) {
		var marchPage MarchPage
		log.Println("Handling at rest routes")

		log.Println("PageURL for page ", r.URL.Path)
		if err := db.One("PageURL", r.URL.Path, &marchPage); err == nil {
			log.Println("PageURL for page ", r.URL.Path)
			renderPage(w, r, marchPage)
		} else {
			renderPage(w, r, MarchPage{PageTemplate: "404.html"})
		}
	})

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	corsRouter := handlers.CORS(headersOk, originsOk, methodsOk)(router)

	if httpErr := http.ListenAndServe(address, corsRouter); httpErr != nil {
		log.Fatalf("Failed to start web service : %s", httpErr.Error())
	} else {
		log.Println("Web service loaded ")
	}

}
