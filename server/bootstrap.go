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

	//MarchPages resource routes
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
			renderPost(w, marchPost)
		} else {
			log.Println("could not fetch route ", r.URL.Path)
		}
	})

	// router.HandleFunc(`/extension`, func(w http.ResponseWriter, r *http.Request) {
	// 	gob.Register(url.Values{})
	// 	switch r.Header.Get("Content-type") {
	// 	case "application/x-www-form-urlencoded":
	// 		{
	// 			resp := new(Response)
	// 			if err := r.ParseForm(); err == nil {
	// 				if err := extensions[r.FormValue("extname")].Call(r.FormValue("extmethod"),
	// 					Request{Input: map[string]interface{}{"form": r.Form},
	// 						Type: "HTML"},
	// 					resp); err == nil {
	// 					log.Println(resp.Output)
	// 					http.Redirect(w, r, r.FormValue("redirectURL"), 301)
	// 				} else {
	// 					log.Fatalln("extensin failed to handle :", err.Error())
	// 				}
	// 			} else {
	// 				log.Fatalln("failed to parse form ", err.Error())
	// 			}
	// 		}
	// 	case "application/json":
	// 		{
	// 			var requestJSON struct {
	// 				ResponseType    string
	// 				Input           string
	// 				ExtensionName   string
	// 				ExtensionMethod string
	// 			}
	// 			resp := new(Response)
	// 			requrstDecoder := json.NewDecoder(r.Body)
	// 			if err := requrstDecoder.Decode(requestJSON); err == nil {
	// 				extensions[requestJSON.ExtensionName].Call(requestJSON.ExtensionMethod, requestJSON.Input, resp)
	// 				fmt.Fprintln(w, resp)
	// 			} else {
	// 				log.Fatalln("could not decode extension request ", err.Error())
	// 			}
	// 		}
	// 	}

	// })

	// Handling pages
	router.HandleFunc(`/{rest:[a-zA-Z0-9=\-\/]*}`, func(w http.ResponseWriter, r *http.Request) {
		var marchPage MarchPage
		log.Println("Handling at rest routes")

		log.Println("PageURL for page ", r.URL.Path)
		if err := db.One("PageURL", r.URL.Path, &marchPage); err == nil {
			log.Println("PageURL for page ", r.URL.Path)
			renderPage(w, marchPage)
		} else {
			renderPage(w, MarchPage{PageTemplate: "404.html"})
		}
	})

	if httpErr := http.ListenAndServe(address, router); httpErr != nil {
		log.Fatalf("Failed to start web service : %s", httpErr.Error())
	} else {
		log.Println("Web service loaded ")
	}

}
