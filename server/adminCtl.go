package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

func adminRoutes(router *mux.Router) {

	router.HandleFunc("/signout", func(w http.ResponseWriter, r *http.Request) {
		session, err := UserSession.Get(r, "mvc-user-session")
		if err == nil {
			for k := range session.Values {
				delete(session.Values, k)
			}
			session.Options.MaxAge = -1
			session.Save(r, w)
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		}
	})

	router.HandleFunc("/login",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				renderAdmin(w, r, "page/login.html", map[string]interface{}{})
			} else if r.Method == "POST" {
				var user MarchUser
				requestBodyDecoder := json.NewDecoder(r.Body)
				if err := requestBodyDecoder.Decode(&user); err != nil {
					log.Println("Failed to decode user data : ", err.Error())
					renderJSON(w, map[string]string{"error": err.Error()})
				} else {
					unhashedPass := user.Password
					log.Println("decoded user data : ", user)
					if err := db.One("Email", user.Email, &user); err == nil {
						user.Password = unhashedPass
						if ok, user := user.LoginUser(); ok {
							issueSession(w, r, user)
							log.Println("session issued for :", user)
							renderJSON(w, map[string]string{"success": "authentication done"})
						} else {
							renderJSON(w, map[string]string{"error": "authentication failed"})
						}
					} else {
						log.Println("could not find user with email:", err)
						renderJSON(w, map[string]string{"error": "email or password is incorrect"})
					}
				}
			} else {
				renderJSON(w, map[string]string{"error": "disallowed http request type"})
			}
		})

	// restricted asset routes
	router.HandleFunc("/admin/themes-thumb/{theme-name}/thumb.png",
		auth(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Returing asset")
			param := mux.Vars(r)
			w.Header().Set("Content-Type", "image/jpeg")
			if thumbFile, err := os.Open(filepath.Join(root, "themes", param["theme-name"], "/thumb.png")); err != nil {
				log.Println("error while eturing asset ", err.Error())
			} else {
				io.Copy(w, thumbFile)
			}
		}))

	// dashboard
	router.HandleFunc("/admin/dashboard",
		auth(func(w http.ResponseWriter, r *http.Request) {
			renderAdmin(w, r, "page/dashboard.html", map[string]interface{}{})
		}))

	// assets
	router.HandleFunc("/admin/assets/{ofType}",
		auth(func(w http.ResponseWriter, r *http.Request) {
			param := mux.Vars(r)
			ofType := param["ofType"]

			dataMap := make(map[string]interface{})

			docMap := make(map[string]interface{})
			imgMap := make(map[string]interface{})
			vidMap := make(map[string]interface{})
			// fetching all files from assets folder
			readFolder := filepath.Join(assetFolder, ofType)
			files, err := ioutil.ReadDir(readFolder)
			if err != nil {
				log.Fatal(err)
			}

			for _, f := range files {
				splits := strings.Split(f.Name(), ".")

				switch splits[len(splits)-1] {
				case "png", "jpg", "gif":
					{
						{
							imgMap[f.Name()] = struct {
								size int64
								url  string
							}{
								f.Size(),
								fmt.Sprintf("/asset/uploaded/images/%s", f.Name()),
							}
						}
					}

				case "mp4", "avi":
					{
						vidMap[f.Name()] = struct {
							size int64
							url  string
						}{
							f.Size(),
							fmt.Sprintf("/asset/uploaded/videos/%s", f.Name()),
						}
					}

				default:
					{
						docMap[f.Name()] = struct {
							size int64
							url  string
						}{
							f.Size(),
							fmt.Sprintf("/asset/uploaded/documents/%s", f.Name()),
						}
					}
				}

			}
			dataMap["docs"] = docMap
			dataMap["imgs"] = imgMap
			dataMap["vids"] = vidMap

			renderAdmin(w, r, "page/assets.html", dataMap)
		}))

	// assets upload handle
	router.HandleFunc("/asset/upload/file",
		auth(func(w http.ResponseWriter, r *http.Request) {
			r.ParseMultipartForm(32 << 20)
			file, handler, err := r.FormFile("file")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()
			fmt.Fprintf(w, "%v", handler.Header)

			splits := strings.Split(handler.Filename, ".")

			storeAt := filepath.Join(assetFolder, "documents", handler.Filename)
			switch splits[len(splits)-1] {
			case "png", "jpg", "gif":
				storeAt = filepath.Join(assetFolder, "images", handler.Filename)
			case "mp4", "avi":
				storeAt = filepath.Join(assetFolder, "videos", handler.Filename)
			default:
			}

			f, err := os.OpenFile(storeAt, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer f.Close()
			io.Copy(f, file)
		}))

	// asset delete handle
	router.HandleFunc("/delete/asset/uploaded/{ofType}/{asset}",
		auth(func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			ofType := vars["ofType"]
			requestedAssetName := vars["asset"]
			if err := os.Remove(filepath.Join(assetFolder, ofType, requestedAssetName)); err == nil {
				renderJSON(w, map[string]string{"error": "requested asset removed"})
			} else {
				renderJSON(w, map[string]string{"error": "Could not delete asset"})
			}
		}))

	// settings
	router.HandleFunc("/admin/settings",
		auth(func(w http.ResponseWriter, r *http.Request) {
			renderAdmin(w, r, "page/settings.html", map[string]interface{}{})
		}))

	router.HandleFunc("/admin/set-theme/{theme-folder-name}",
		auth(func(w http.ResponseWriter, r *http.Request) {
			param := mux.Vars(r)
			config.Theme = param["theme-folder-name"]
			if err := db.Save(&config); err != nil {
				log.Fatalln("Failed to update theme ", err.Error())
			}
			renderAdmin(w, r, "page/settings.html", map[string]interface{}{})
		}))

	// Menu management routes
	router.HandleFunc("/admin/menus/list",
		auth(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				allMenus := []MarchMenu{}
				if err := db.All(&allMenus); err == nil {
					renderAdmin(w, r, "page/menus.html", map[string]interface{}{
						"menus": allMenus,
					})
				}
			}
		}))

	router.HandleFunc("/admin/menu/create",
		auth(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				menus := new([]MarchMenu)
				if err := db.All(menus); err != nil {
					renderJSON(w, map[string]interface{}{
						"error": fmt.Sprint("Could not get all menus ", err.Error()),
					})
					return
				}
				renderAdmin(w, r, "page/menu-editor.html", map[string]interface{}{
					"menus": menus,
				})
			}
		}))

	router.HandleFunc("/admin/menu/save",
		auth(func(w http.ResponseWriter, r *http.Request) {
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

				// Setting authenticated user as menu creator
				usession, _ := UserSession.Get(r, "mvc-user-session")
				menu.MarchUserID = usession.Values["id"].(int)
				menu.Co = time.Now()
				if err := db.Save(&menu); err != nil {
					log.Println("Faild to save menu :", err.Error())
					renderJSON(w, map[string]interface{}{
						"Error": err.Error(),
					})
				} else {
					log.Println("Menu saved succesfully - ", menu)
					renderJSON(w, map[string]interface{}{
						"Succcess": true,
					})
				}
			}
		}))

	router.HandleFunc("/admin/site/settings",
		auth(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
					renderJSON(w, map[string]Config{"config": config})
					return
				}
			} else if r.Method == "POST" {
				var inputStruct struct {
					Name       string `json:"Name"`
					LogoURL    string `json:"LogoURL"`
					FaviconURL string `json:"FaviconURL"`
				}
				// decoding input
				inputDecoder := json.NewDecoder(r.Body)
				if err := inputDecoder.Decode(&inputStruct); err == nil {

					if err := db.One("Status", "Active", &config); err == nil {
						config.Name = inputStruct.Name
						config.FaviconURL = inputStruct.FaviconURL
						config.LogoURL = inputStruct.LogoURL
						if err := db.Save(&config); err != nil {
							log.Fatalln("failed to save config", err.Error())
						}
					} else {
						log.Fatalln("Could not save local config", err.Error())
					}
				} else {
					log.Fatalln("Failed to decode input", err.Error())
				}
			}

		}))

	router.HandleFunc("/admin/theme/settings",
		auth(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
					renderJSON(w, map[string]ThemeConfig{"config": themeConfig})
					return
				}
			} else if r.Method == "POST" {
				inputDecoder := json.NewDecoder(r.Body)
				var menus = struct {
					themeMenus []*ThemeMenu `json:menus`
				}{}

				if err := inputDecoder.Decode(&menus); err != nil {
					renderJSON(w, map[string]string{"error": err.Error()})
					return
				}
				log.Println(menus.themeMenus)
				// if err := db.UpdateField(&ThemeConfig{ID: themeConfig.ID}, "Menus", menus.themeMenus); err != nil {
				// 	log.Println("failed to update menu settings ", err.Error())
				// 	renderJSON(w, map[string]string{"error": err.Error()})
				// 	return
				// }

				// themeConfig.Menus = menus.themeMenus

				renderJSON(w, map[string]interface{}{"success": themeConfig.Menus})
			}

		}))

	router.HandleFunc("/admin/menu/{ID}/edit",
		auth(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				vars := mux.Vars(r)
				if _menuID, _conErr := strconv.Atoi(vars["ID"]); _conErr == nil {
					if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
						_menu := MarchMenu{}
						if err := db.One("ID", _menuID, &_menu); err == nil {
							renderJSON(w, map[string]MarchMenu{
								"menu": _menu,
							})
						} else {
							renderJSON(w, map[string]string{
								"error": err.Error(),
							})
						}
					} else {
						renderAdmin(w, r, "page/menu-editor.html", map[string]interface{}{
							"ID": _menuID,
						})
					}
				} else {
					fmt.Fprintln(w, "Wrong URL ", _conErr.Error())
				}
			}
		}))

	// Page managemnt routes
	router.HandleFunc("/admin/page/create",
		auth(author(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				renderAdmin(w, r, "page/page-create.html", map[string]interface{}{})
			} else if r.Method == "POST" {

				var requestData struct {
					Desc, HTML, Keywords, PageTitle, PageURL, PageTemplate string
				}
				requestDecoder := json.NewDecoder(r.Body)
				requestDecoder.Decode(&requestData)

				if uuid, err := uuid.NewV1(); err != nil {
					log.Fatalln("Failed to gerate page id :", err.Error())
				} else {
					usession, _ := UserSession.Get(r, "mvc-user-session")
					newPage := MarchPage{
						MarchUserID:  usession.Values["id"].(int),
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
					// newPage.Uo = time.Now()
					if err := db.Save(&newPage); err != nil {
						log.Fatalln("failed to save page : ", err.Error())
					} else {
						log.Println("Route & page saved with id ", newPage)
						renderJSON(w, newPage)
					}
				}
			}
		})))

	router.HandleFunc("/admin/page/{id}/edit",
		auth(func(w http.ResponseWriter, r *http.Request) {
			param := mux.Vars(r)
			MarchPage := MarchPage{}
			if err := db.One("PageNumber", param["id"], &MarchPage); err == nil {
				if r.Method == "GET" {
					renderAdmin(w, r, "page/page-edit.html", map[string]interface{}{
						"page": MarchPage,
					})
				} else {
					var requestData struct {
						Desc, HTML, Keywords, PageTitle, PageURL, PageTemplate string
					}
					if !originalWriter(MarchPage.MarchUserID, r) {
						log.Println(http.StatusUnauthorized, " page editing blocked.")
						return
					}
					requestDecoder := json.NewDecoder(r.Body)
					requestDecoder.Decode(&requestData)
					usession, _ := UserSession.Get(r, "mvc-user-session")
					MarchPage.UpdaterID = usession.Values["id"].(int)
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

		}))

	router.HandleFunc("/admin/page/{id}/delete",
		auth(func(w http.ResponseWriter, r *http.Request) {
			param := mux.Vars(r)
			MarchPage := MarchPage{}
			if err := db.One("PageNumber", param["id"], &MarchPage); err == nil {
				usession, _ := UserSession.Get(r, "mvc-user-session")
				MarchPage.UpdaterID = usession.Values["id"].(int)
				MarchPage.Do = time.Now()
				if err := db.DeleteStruct(&MarchPage); err == nil {
					http.Redirect(w, r, "/admin/pages/list", 301)
				} else {
					http.Redirect(w, r, "/admin/pages/list", 301)
				}
			} else {
				log.Fatalln("couldn not get page for corrsponding route")
			}

		}))

	router.HandleFunc("/admin/pages/list",
		auth(func(w http.ResponseWriter, r *http.Request) {
			var pages []MarchPage
			if err := db.All(&pages); err != nil {
				log.Fatalln("failed to load routes : ", err.Error())
			} else {
				renderAdmin(w, r, "page/pages.html", map[string]interface{}{
					"pages": pages,
				})
			}
		}))

	// Post managemnt routes
	router.HandleFunc("/admin/post/create",
		auth(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "GET" {
				renderAdmin(w, r, "page/post-create.html", map[string]interface{}{})
			} else if r.Method == "POST" {

				var requestData struct {
					Desc, HTML, Keywords, PageTitle, PageURL, PageTemplate, Tag1, Tag2, Tag3 string
				}
				requestDecoder := json.NewDecoder(r.Body)
				requestDecoder.Decode(&requestData)

				if uuid, err := uuid.NewV1(); err != nil {
					log.Fatalln("Failed to gerate page id :", err.Error())
				} else {
					usession, _ := UserSession.Get(r, "mvc-user-session")
					newpost := MarchPost{
						MarchUserID:  usession.Values["id"].(int),
						PageTemplate: requestData.PageTemplate,
						PageURL:      requestData.PageURL,
						PageTitle:    requestData.PageTitle,
						Tag1:         requestData.Tag1,
						Tag2:         requestData.Tag2,
						Tag3:         requestData.Tag3,
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

		}))

	router.HandleFunc("/admin/post/{id}/edit",
		auth(func(w http.ResponseWriter, r *http.Request) {
			param := mux.Vars(r)

			MarchPost := MarchPost{}
			if err := db.One("PageNumber", param["id"], &MarchPost); err == nil {

				if r.Method == "GET" {
					renderAdmin(w, r, "page/post-edit.html", map[string]interface{}{
						"post": MarchPost,
					})
					return
				} else {
					if !originalWriter(MarchPost.MarchUserID, r) {
						log.Println(http.StatusUnauthorized, " post editing blocked.")
						return
					}
					var requestData struct {
						Desc, HTML, Keywords, PageTitle, PageURL, PageTemplate, Tag1, Tag2, Tag3 string
					}
					requestDecoder := json.NewDecoder(r.Body)
					requestDecoder.Decode(&requestData)
					usession, _ := UserSession.Get(r, "mvc-user-session")
					MarchPost.UpdaterID = usession.Values["id"].(int)
					MarchPost.PageTemplate = requestData.PageTemplate
					MarchPost.Tag1 = requestData.Tag1
					MarchPost.Tag2 = requestData.Tag2
					MarchPost.Tag3 = requestData.Tag3
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

		}))

	router.HandleFunc("/admin/post/{id}/delete",
		auth(func(w http.ResponseWriter, r *http.Request) {
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

		}))

	router.HandleFunc("/admin/posts/list",

		auth(func(w http.ResponseWriter, r *http.Request) {
			var posts []MarchPost
			if err := db.All(&posts); err != nil {
				log.Fatalln("failed to load routes : ", err.Error())
			} else {
				renderAdmin(w, r, "page/posts.html", map[string]interface{}{
					"posts": posts,
				})
			}
		}))

	// User management routes
	router.HandleFunc("/admin/users/list",
		auth(func(w http.ResponseWriter, r *http.Request) {
			var users []MarchUser
			if err := db.All(&users); err != nil {
				log.Fatalln("failed to load routes : ", err.Error())
			} else {
				renderAdmin(w, r, "page/users.html", map[string]interface{}{
					"users": users,
				})
			}
		}))
	// main route

	router.HandleFunc("/admin",
		auth(func(w http.ResponseWriter, r *http.Request) {

			renderAdmin(w, r, "page/index.html", map[string]interface{}{})

		}))

}
