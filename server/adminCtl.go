package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

func adminRoutes(router *mux.Router) {

	// dashboard
	router.HandleFunc("/admin/dashboard", func(w http.ResponseWriter, r *http.Request) {
		renderAdmin(w, "page/dashboard.html", map[string]interface{}{})
	})

	// content routes
	router.HandleFunc("/admin/page/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			renderAdmin(w, "page/page-create.html", map[string]interface{}{})
		} else if r.Method == "POST" {

			var requestData struct {
				Desc, HTML, Keywords, PageTitle, PageURL string
			}
			requestDecoder := json.NewDecoder(r.Body)
			requestDecoder.Decode(&requestData)
			newRoute := SlingRoute{
				PageURL: requestData.PageURL,
			}
			if uuid, err := uuid.NewV1(); err != nil {
				log.Fatalln("Failed to gerate page id :", err.Error())
			} else {
				newRoute.PageNumber = uuid.String()
				if err := db.Save(&newRoute); err != nil {
					log.Fatalln("Failed to save route :", err.Error())
				} else {
					newPage := SlingPage{
						PageTitle: requestData.PageTitle,
						Content: SlingPageContent{
							Descritpion: requestData.Desc,
							HTML:        requestData.HTML,
							Keywords:    requestData.Keywords,
						},
					}
					newPage.PageNumber = newRoute.PageNumber
					if err := db.Save(&newPage); err != nil {
						log.Fatalln("failed to save page : ", err.Error())
					} else {
						log.Println("Route & page saved with id ", newRoute.PageNumber, newRoute.PageURL)
						fmt.Fprintln(w, newPage)
					}
				}
			}
		}
	})

	router.HandleFunc("/admin/page/{id}/edit", func(w http.ResponseWriter, r *http.Request) {
		renderAdmin(w, "page/page-edit.html", map[string]interface{}{})
	})

	router.HandleFunc("/admin/pages/list", func(w http.ResponseWriter, r *http.Request) {
		var pages []SlingPage
		if err := db.All(&pages); err != nil {
			log.Fatalln("failed to load routes : ", err.Error())
		} else {
			renderAdmin(w, "page/pages.html", map[string]interface{}{
				"pages": pages,
			})
		}
	})

	// settigns routes

	// main route
	router.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		renderAdmin(w, "page/index.html", map[string]interface{}{})
	})

}
