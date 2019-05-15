package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

func adminRoutes(router *mux.Router) {

	// restricted asset routes
	router.HandleFunc("/admin/themes-thumb/{theme-name}/thumb.png",
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("Returing asset")
			param := mux.Vars(r)
			w.Header().Set("Content-Type", "image/jpeg")
			if thumbFile, err := os.Open(filepath.Join(root, "themes", param["theme-name"], "/thumb.png")); err != nil {
				log.Println("error while eturing asset ", err.Error())
			} else {
				io.Copy(w, thumbFile)
			}
		})

	// extensions routes
	router.HandleFunc("/admin/extensions/{extension-name}",
		func(w http.ResponseWriter, r *http.Request) {
			log.Println("Starting Appliction")
			param := mux.Vars(r)
			extensionPath := filepath.Join(extensionFolder, param["extension-name"], "hello.exe")
			cmd := exec.Command(extensionPath)
			if err := cmd.Start(); err == nil {
				loadExtensions()
				// http.Redirect(w, r, "/admin/settings", 301)
			} else {
				log.Fatal(err.Error())
			}
		})

	// dashboard
	router.HandleFunc("/admin/dashboard", func(w http.ResponseWriter, r *http.Request) {
		renderAdmin(w, "page/dashboard.html", map[string]interface{}{})
	})

	// Settings Route
	router.HandleFunc("/admin/settings", func(w http.ResponseWriter, r *http.Request) {
		renderAdmin(w, "page/settings.html", map[string]interface{}{})
	})

	router.HandleFunc("/admin/set-theme/{theme-folder-name}", func(w http.ResponseWriter, r *http.Request) {
		param := mux.Vars(r)
		config.Theme = param["theme-folder-name"]
		if err := db.Save(&config); err != nil {
			log.Fatalln("Failed to update theme ", err.Error())
		}
		renderAdmin(w, "page/settings.html", map[string]interface{}{})
	})

	// Page managemnt routes
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

	// Post managemnt routes
	router.HandleFunc("/admin/post/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			renderAdmin(w, "page/post-create.html", map[string]interface{}{})
		} else if r.Method == "POST" {

			var requestData struct {
				Desc, HTML, Keywords, PageTitle, PageURL, PageTemplate string
			}
			requestDecoder := json.NewDecoder(r.Body)
			requestDecoder.Decode(&requestData)

			if uuid, err := uuid.NewV1(); err != nil {
				log.Fatalln("Failed to gerate page id :", err.Error())
			} else {

				newpost := SlingPost{
					PageTemplate: requestData.PageTemplate,
					PageURL:      requestData.PageURL,
					PageTitle:    requestData.PageTitle,
					Content: SlingPageContent{
						Desc:     requestData.Desc,
						HTML:     requestData.HTML,
						Keywords: requestData.Keywords,
					},
				}
				newpost.PageNumber = uuid.String()
				newpost.Co = time.Now()
				newpost.Uo = time.Now()
				if err := db.Save(&newpost); err != nil {
					log.Fatalln("failed to save page : ", err.Error())
				} else {
					log.Println("Route & page saved with id ", newpost)
					renderJSON(w, newpost)
				}
			}
		}

	})

	router.HandleFunc("/admin/post/{id}/edit", func(w http.ResponseWriter, r *http.Request) {
		param := mux.Vars(r)

		slingpost := SlingPost{}
		if err := db.One("PageNumber", param["id"], &slingpost); err == nil {

			if r.Method == "GET" {
				renderAdmin(w, "page/post-edit.html", map[string]interface{}{
					"post": slingpost,
				})
				return
			} else {

				var requestData struct {
					Desc, HTML, Keywords, PageTitle, PageURL, PageTemplate string
				}
				requestDecoder := json.NewDecoder(r.Body)
				requestDecoder.Decode(&requestData)
				slingpost.PageTemplate = requestData.PageTemplate
				slingpost.PageURL = requestData.PageURL
				slingpost.PageNumber = param["id"]
				slingpost.PageTitle = requestData.PageTitle
				slingpost.Content.Desc = requestData.Desc
				slingpost.Content.HTML = requestData.HTML
				slingpost.Content.Keywords = requestData.Keywords
				slingpost.Uo = time.Now()

				db.Save(&slingpost)

				return
			}

		} else {
			log.Fatalln("couldn not get page for corrsponding route")
		}

	})

	router.HandleFunc("/admin/post/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		param := mux.Vars(r)

		slingpost := SlingPost{}
		if err := db.One("PageNumber", param["id"], &slingpost); err == nil {
			if err := db.DeleteStruct(&slingpost); err == nil {
				http.Redirect(w, r, "/admin/posts/list", 301)
			} else {
				http.Redirect(w, r, "/admin/posts/list", 301)
			}
		} else {
			log.Fatalln("couldn not get page for corrsponding route")
		}

	})

	router.HandleFunc("/admin/posts/list", func(w http.ResponseWriter, r *http.Request) {
		var posts []SlingPost
		if err := db.All(&posts); err != nil {
			log.Fatalln("failed to load routes : ", err.Error())
		} else {
			renderAdmin(w, "page/posts.html", map[string]interface{}{
				"posts": posts,
			})
		}
	})

	// main route
	router.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		renderAdmin(w, "page/index.html", map[string]interface{}{})
	})

}
