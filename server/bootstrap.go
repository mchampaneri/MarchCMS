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

	log.Println("Listening on ", address)
	// start web on spacific address
	router := mux.NewRouter()

	router.PathPrefix("/sl-res/").
		Handler(http.StripPrefix("/sl-res/",
			http.FileServer(http.Dir("./admin/resource"))))

	router.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		renderAdmin(w, "page/index.html", map[string]interface{}{})
	})

	var routes []SlingRoute
	if err := db.All(&routes); err != nil {
		log.Fatalln("failed to load routes : ", err.Error())
	} else {
		for _, route := range routes {
			mountRoute(route, router)
		}
	}

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
