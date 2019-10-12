// reading and loading configuration
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func systemCheck() {

	// check if theme folder is existing.
	if _, err := os.Open(themesFolder); err != nil {
		log.Println("Can not find theme folder :", err.Error())
		if err := os.Mkdir(themesFolder, os.ModePerm); err != nil {
			log.Println("Can not create new theme folder :", err.Error())
		}
	}

	// check if admin folder is existing.
	if _, err := os.Open(adminFolder); err != nil {
		log.Println("Can not find admin folder :", err.Error())
		if err := os.Mkdir(adminFolder, os.ModePerm); err != nil {
			log.Println("Can not create new admin folder :", err.Error())
		}
	}

	// check if assets folder is existing.
	if _, err := os.Open(assetFolder); err != nil {
		log.Println("Can not find asset folder :", err.Error())
		if err := os.Mkdir(assetFolder, os.ModePerm); err != nil {
			log.Println("Can not create new asset folder :", err.Error())
		}
		if err := os.Mkdir(filepath.Join(assetFolder, "images"), os.ModePerm); err != nil {
			log.Println("Can not create new asset folder :", err.Error())
		}
		if err := os.Mkdir(filepath.Join(assetFolder, "videos"), os.ModePerm); err != nil {
			log.Println("Can not create new asset folder :", err.Error())
		}
		if err := os.Mkdir(filepath.Join(assetFolder, "documents"), os.ModePerm); err != nil {
			log.Println("Can not create new asset folder :", err.Error())
		}
	}
}

func loadSiteConfig(config *Config) {
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

func loadThemeConfig(theme string, config *ThemeConfig) {
	// Read and load configurtion json
	if file, fileErr := os.Open(path.Join(themesFolder, theme, "./config.json")); fileErr != nil {
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

			marchUser := MarchUser{}
			if err := db.One("ID", marchPost.MarchUserID, &marchUser); err == nil {
				marchPost.MarchUserObj = marchUser
			}

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
			marchUser := MarchUser{}
			if err := db.One("ID", marchPage.MarchUserID, &marchUser); err == nil {
				marchPage.MarchUserObj = marchUser
			}
			renderPage(w, r, marchPage)
		} else {
			renderPage(w, r, MarchPage{PageTemplate: "404.html"})
		}
	})

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	corsRouter := handlers.CORS(headersOk, originsOk, methodsOk)(router)

	if httpErr := http.ListenAndServe(address, statusMiddleware(corsRouter)); httpErr != nil {
		log.Fatalf("Failed to start web service : %s", httpErr.Error())
	} else {
		log.Println("Web service loaded ")
	}

}

func statusMiddleware(next http.Handler) http.Handler {
	if config.Live == "Yes" {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Our middleware logic goes here...
			next.ServeHTTP(w, r)
		})
	} else {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Our middleware logic goes here...
			fmt.Fprintln(w, "Site under maintianance")
		})
	}
}
