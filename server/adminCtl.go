package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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

	// Menu management routes
	router.HandleFunc("/admin/menus/list", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			renderAdmin(w, "page/menus.html", map[string]interface{}{})
		}
	})

	router.HandleFunc("/admin/menu/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			menus := new([]MarchMenu)
			if err := db.All(menus); err != nil {
				renderJSON(w, map[string]interface{}{
					"error": fmt.Sprint("Could not get all menus ", err.Error()),
				})
				return
			}
			renderAdmin(w, "page/menu-create.html", map[string]interface{}{
				"menus": menus,
			})
		}
	})

	router.HandleFunc("/admin/menu/save", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			menu := MarchMenu{}
			inputDecoder := json.NewDecoder(r.Body)
			defer r.Body.Close()

			// decoding menu input
			if err := inputDecoder.Decode(&menu); err != nil {
				renderJSON(w, map[string]interface{}{
					"error": fmt.Sprint("Error during decoding input :", err.Error()),
				})
				return
			}
			// Sluggifying menu title and menuItme title
			menu.Slug = Slugy([]string{menu.Name})
			for _, menuItem := range menu.Items {
				menuItem.Item.Slug = Slugy([]string{menuItem.Item.Title})
			}
			if err := db.Save(&menu); err != nil {
				log.Fatal("Faild to save menu :", err.Error())
			} else {
				renderJSON(w, map[string]interface{}{
					"Succcess": true,
				},
				)
			}
		}
	})

	router.HandleFunc("/admin/menu/edit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			renderAdmin(w, "page/menu.html", map[string]interface{}{})
		}
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

				newPage := MarchPage{
					PageTemplate: requestData.PageTemplate,
					PageURL:      requestData.PageURL,
					PageTitle:    requestData.PageTitle,
					Content: MarchPageContent{
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

		MarchPage := MarchPage{}
		if err := db.One("PageNumber", param["id"], &MarchPage); err == nil {

			if r.Method == "GET" {
				renderAdmin(w, "page/page-edit.html", map[string]interface{}{
					"page": MarchPage,
				})
				return
			} else {

				var requestData struct {
					Desc, HTML, Keywords, PageTitle, PageURL, PageTemplate string
				}
				requestDecoder := json.NewDecoder(r.Body)
				requestDecoder.Decode(&requestData)
				MarchPage.PageTemplate = requestData.PageTemplate
				MarchPage.PageURL = requestData.PageURL
				MarchPage.PageNumber = param["id"]
				MarchPage.PageTitle = requestData.PageTitle
				MarchPage.Content.Desc = requestData.Desc
				MarchPage.Content.HTML = requestData.HTML
				MarchPage.Content.Keywords = requestData.Keywords
				MarchPage.Uo = time.Now()
				db.Save(&MarchPage)
				return
			}

		} else {
			log.Fatalln("couldn not get page for corrsponding route")
		}

	})

	router.HandleFunc("/admin/page/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		param := mux.Vars(r)
		MarchPage := MarchPage{}
		if err := db.One("PageNumber", param["id"], &MarchPage); err == nil {
			if err := db.DeleteStruct(&MarchPage); err == nil {
				http.Redirect(w, r, "/admin/pages/list", 301)
			} else {
				http.Redirect(w, r, "/admin/pages/list", 301)
			}
		} else {
			log.Fatalln("couldn not get page for corrsponding route")
		}

	})

	router.HandleFunc("/admin/pages/list", func(w http.ResponseWriter, r *http.Request) {
		var pages []MarchPage
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

				newpost := MarchPost{
					PageTemplate: requestData.PageTemplate,
					PageURL:      requestData.PageURL,
					PageTitle:    requestData.PageTitle,
					Content: MarchPageContent{
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

		MarchPost := MarchPost{}
		if err := db.One("PageNumber", param["id"], &MarchPost); err == nil {

			if r.Method == "GET" {
				renderAdmin(w, "page/post-edit.html", map[string]interface{}{
					"post": MarchPost,
				})
				return
			} else {

				var requestData struct {
					Desc, HTML, Keywords, PageTitle, PageURL, PageTemplate string
				}
				requestDecoder := json.NewDecoder(r.Body)
				requestDecoder.Decode(&requestData)
				MarchPost.PageTemplate = requestData.PageTemplate
				MarchPost.PageURL = requestData.PageURL
				MarchPost.PageNumber = param["id"]
				MarchPost.PageTitle = requestData.PageTitle
				MarchPost.Content.Desc = requestData.Desc
				MarchPost.Content.HTML = requestData.HTML
				MarchPost.Content.Keywords = requestData.Keywords
				MarchPost.Uo = time.Now()

				db.Save(&MarchPost)

				return
			}

		} else {
			log.Fatalln("couldn not get page for corrsponding route")
		}

	})

	router.HandleFunc("/admin/post/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		param := mux.Vars(r)

		MarchPost := MarchPost{}
		if err := db.One("PageNumber", param["id"], &MarchPost); err == nil {
			if err := db.DeleteStruct(&MarchPost); err == nil {
				http.Redirect(w, r, "/admin/posts/list", 301)
			} else {
				http.Redirect(w, r, "/admin/posts/list", 301)
			}
		} else {
			log.Fatalln("couldn not get page for corrsponding route")
		}

	})

	router.HandleFunc("/admin/posts/list", func(w http.ResponseWriter, r *http.Request) {
		var posts []MarchPost
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
