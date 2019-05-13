package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

func adminRoutes(router *mux.Router) {

	// dashboard
	router.HandleFunc("/admin/dashboard", func(w http.ResponseWriter, r *http.Request) {
		renderAdmin(w, "page/dashboard.html", map[string]interface{}{})
	})

	// content routes
	router.HandleFunc("/admin/page/create-v1", func(w http.ResponseWriter, r *http.Request) {
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
							Desc:     requestData.Desc,
							HTML:     requestData.HTML,
							Keywords: requestData.Keywords,
						},
					}
					newPage.PageNumber = newRoute.PageNumber
					if err := db.Save(&newPage); err != nil {
						log.Fatalln("failed to save page : ", err.Error())
					} else {
						log.Println("Route & page saved with id ", newRoute.PageNumber, newRoute.PageURL)
						renderJSON(w, newPage)
					}
				}
			}
		}
	})

	router.HandleFunc("/admin/page/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			renderAdmin(w, "page/page-create.html", map[string]interface{}{})
		} else if r.Method == "POST" {

			var requestData struct {
				Desc, HTML, Keywords, PageTitle, PageURL, PageTemplate string
			}
			requestDecoder := json.NewDecoder(r.Body)
			requestDecoder.Decode(&requestData)

			if uuid, err := uuid.NewV1(); err != nil {
				log.Fatalln("Failed to gerate page id :", err.Error())
			} else {

				newPage := SlingPage{
					PageTemplate: requestData.PageTemplate,
					PageURL:      requestData.PageURL,
					PageTitle:    requestData.PageTitle,
					Content: SlingPageContent{
						Desc:     requestData.Desc,
						HTML:     requestData.HTML,
						Keywords: requestData.Keywords,
					},
				}
				newPage.PageNumber = uuid.String()
				newPage.Co = time.Now()
				newPage.Uo = time.Now()
				if err := db.Save(&newPage); err != nil {
					log.Fatalln("failed to save page : ", err.Error())
				} else {
					log.Println("Route & page saved with id ", newPage)
					renderJSON(w, newPage)
				}
			}
		}

	})

	router.HandleFunc("/admin/page/{id}/edit-v1", func(w http.ResponseWriter, r *http.Request) {
		param := mux.Vars(r)
		slingRoute := SlingRoute{}

		if err := db.One("PageNumber", param["id"], &slingRoute); err == nil {
			slingPage := SlingPage{}
			if err := db.One("PageNumber", param["id"], &slingPage); err == nil {

				if r.Method == "GET" {
					renderAdmin(w, "page/page-edit.html", map[string]interface{}{
						"route": slingRoute,
						"page":  slingPage,
					})
					return
				} else {

					var requestData struct {
						Desc, HTML, Keywords, PageTitle, PageURL string
					}
					requestDecoder := json.NewDecoder(r.Body)
					requestDecoder.Decode(&requestData)

					slingRoute.PageURL = requestData.PageURL
					slingRoute.PageNumber = param["id"]

					slingPage.PageNumber = param["id"]
					slingPage.PageTitle = requestData.PageTitle
					slingPage.Content.Desc = requestData.Desc
					slingPage.Content.HTML = requestData.HTML
					slingPage.Content.Keywords = requestData.Keywords
					db.Save(&slingPage)
					db.Save(&slingRoute)
					return
				}

			} else {
				log.Fatalln("couldn not get page for corrsponding route")
			}
		} else {
			log.Fatalln("couldn not get route for corrsponding route")
		}

	})

	router.HandleFunc("/admin/page/{id}/edit", func(w http.ResponseWriter, r *http.Request) {
		param := mux.Vars(r)

		slingPage := SlingPage{}
		if err := db.One("PageNumber", param["id"], &slingPage); err == nil {

			if r.Method == "GET" {
				renderAdmin(w, "page/page-edit.html", map[string]interface{}{
					"page": slingPage,
				})
				return
			} else {

				var requestData struct {
					Desc, HTML, Keywords, PageTitle, PageURL, PageTemplate string
				}
				requestDecoder := json.NewDecoder(r.Body)
				requestDecoder.Decode(&requestData)
				slingPage.PageTemplate = requestData.PageTemplate
				slingPage.PageURL = requestData.PageURL
				slingPage.PageNumber = param["id"]
				slingPage.PageTitle = requestData.PageTitle
				slingPage.Content.Desc = requestData.Desc
				slingPage.Content.HTML = requestData.HTML
				slingPage.Content.Keywords = requestData.Keywords
				slingPage.Uo = time.Now()

				db.Save(&slingPage)

				return
			}

		} else {
			log.Fatalln("couldn not get page for corrsponding route")
		}

	})

	router.HandleFunc("/admin/page/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		param := mux.Vars(r)

		slingPage := SlingPage{}
		if err := db.One("PageNumber", param["id"], &slingPage); err == nil {
			if err := db.DeleteStruct(&slingPage); err == nil {
				http.Redirect(w, r, "/admin/pages/list", 301)
			} else {
				http.Redirect(w, r, "/admin/pages/list", 301)
			}
		} else {
			log.Fatalln("couldn not get page for corrsponding route")
		}

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
