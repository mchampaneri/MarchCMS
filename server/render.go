// generatign responses in various foramts
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"reflect"

	"github.com/CloudyKit/jet"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

var frontInstance = jet.NewHTMLSet(themesFolder)
var adminInstance = jet.NewHTMLSet(adminFolder)

func init() {

	frontInstance.SetDevelopmentMode(true)
	adminInstance.SetDevelopmentMode(true)

	/////////////// Front ///////////////

	frontInstance.AddGlobalFunc("SiteAddress", func(a jet.Arguments) reflect.Value {
		return reflect.ValueOf(config.Address)
	})

	frontInstance.AddGlobalFunc("SiteTitle", func(a jet.Arguments) reflect.Value {
		return reflect.ValueOf(config.Name)
	})

	frontInstance.AddGlobalFunc("PageList", func(a jet.Arguments) reflect.Value {
		var AllPages []MarchPage
		db.All(&AllPages)
		return reflect.ValueOf(AllPages)
	})

	frontInstance.AddGlobalFunc("Menu", func(a jet.Arguments) reflect.Value {
		_menu := MarchMenu{}
		input := a.Get(0).String()

		for _, menu := range themeConfig.Menus {
			if menu.Place == input {
				log.Println("we found menu ", menu.Place)
				if err := db.One("Slug", menu.Menu, &_menu); err == nil {
					return reflect.ValueOf(_menu)
				} else {
					return reflect.ValueOf(fmt.Sprint(menu.Place, " menu missig"))
				}

			}
		}

		return reflect.ValueOf("-")

	})

	frontInstance.AddGlobalFunc("PostByTag", func(a jet.Arguments) reflect.Value {
		// wg := &sync.WaitGroup{}
		input := a.Get(0).String()
		// wg.Add(3)
		var postsTag1, postsTag2, postsTag3 []MarchPost

		finalList := make(map[string]MarchPost, 1000)

		// go func() {
		if err := db.Find("Tag1", input, &postsTag1); err == nil {
			for _, post := range postsTag1 {
				finalList[post.PageNumber] = post
			}
		}
		log.Println("Tag1 search complete.")
		// wg.Done()
		// }()

		// go func() {
		if err := db.Find("Tag2", input, &postsTag2); err == nil {
			for _, post := range postsTag2 {
				finalList[post.PageNumber] = post
			}
		}
		log.Println("Tag2 search complete.")
		// wg.Done()
		// }()

		// go func() {
		if err := db.Find("Tag3", input, &postsTag3); err == nil {
			for _, post := range postsTag3 {
				finalList[post.PageNumber] = post
			}
		}
		log.Println("Tag3 search complete.")
		// wg.Done()
		// }()

		// wg.Wait()
		return reflect.ValueOf(finalList)

	})
	////////////// Admin //////////////////
	adminInstance.AddGlobalFunc("SiteAddress", func(a jet.Arguments) reflect.Value {
		return reflect.ValueOf(config.Address)
	})

	adminInstance.AddGlobalFunc("SiteTitle", func(a jet.Arguments) reflect.Value {
		return reflect.ValueOf(config.Name)
	})

	adminInstance.AddGlobalFunc("ActiveTheme", func(a jet.Arguments) reflect.Value {
		return reflect.ValueOf(config.Theme)
	})

	adminInstance.AddGlobalFunc("PageTemplates", func(a jet.Arguments) reflect.Value {
		// Read Pages templates
		templates := make([]string, 0, 10)
		if fileInfo, err := ioutil.ReadDir(filepath.Join(root, "themes", config.Theme, "pages")); err == nil {
			for _, file := range fileInfo {
				templates = append(templates, file.Name())
			}
			return reflect.ValueOf(templates)
		}
		return reflect.ValueOf(nil)
	})

	adminInstance.AddGlobalFunc("PostTemplates", func(a jet.Arguments) reflect.Value {
		// Read Pages templates
		templates := make([]string, 0, 10)
		if fileInfo, err := ioutil.ReadDir(filepath.Join(root, "themes", config.Theme, "posts")); err == nil {
			for _, file := range fileInfo {
				templates = append(templates, file.Name())
			}
			return reflect.ValueOf(templates)
		}
		return reflect.ValueOf(nil)
	})

	adminInstance.AddGlobalFunc("InstalledThemes", func(a jet.Arguments) reflect.Value {
		// Read Pages templates
		templates := make([]struct{ Name, Thumb string }, 0, 10)
		if fileInfo, err := ioutil.ReadDir(filepath.Join(root, "themes")); err == nil {
			for _, file := range fileInfo {
				templates = append(templates, struct{ Name, Thumb string }{file.Name(),
					fmt.Sprint("/admin/themes-thumb/", file.Name(), "/thumb.png")})
			}
			return reflect.ValueOf(templates)
		}
		return reflect.ValueOf(nil)
	})

}

func renderPage(w io.Writer, r *http.Request, page MarchPage) {
	var pageTemplate = "index.html"
	if page.PageTemplate != "" && page.PageTemplate != "-" {
		pageTemplate = filepath.Join(config.Theme, "pages", page.PageTemplate)
	}
	t, err := frontInstance.GetTemplate(pageTemplate)
	if err != nil {
		log.Println(pageTemplate, " - ", config.Theme, " - ", err.Error())
		return
	}

	dataMap := map[string]interface{}{
		"Page": page,
	}
	log.Println(page)

	output := blackfriday.Run([]byte(page.Content.HTML))
	dataMap["output"] = output
	dataMap["requestURL"] = r.RequestURI
	if err = t.Execute(w, nil, dataMap); err != nil {
		log.Println(" render.go  View  : %s", err.Error())
	}
}

func renderPost(w io.Writer, r *http.Request, post MarchPost) {
	var pageTemplate = "index.html"
	if post.PageTemplate != "" && post.PageTemplate != "-" {
		pageTemplate = filepath.Join(config.Theme, "posts", post.PageTemplate)
	}
	t, err := frontInstance.GetTemplate(pageTemplate)
	if err != nil {
		log.Println(pageTemplate, " - ", config.Theme, " - ", err.Error())
	}
	dataMap := map[string]interface{}{
		"Page": post,
	}
	output := blackfriday.Run([]byte(post.Content.HTML))
	dataMap["output"] = output
	dataMap["requestURL"] = r.RequestURI
	if err = t.Execute(w, nil, dataMap); err != nil {
		log.Println(" - respnose-generator.go  View  : %s", err.Error())
	}
}

func renderAdmin(w io.Writer, r *http.Request, page string, dataMap map[string]interface{}) {
	// log.Println("Render admin is executing")
	// log.Fatalln(page)
	if t, err := adminInstance.GetTemplate(page); err == nil {
		// dataMap := map[string]interface{}{}s
		dataMap["requestURL"] = r.RequestURI
		usession, _ := UserSession.Get(r, "mvc-user-session")
		dataMap["authUser"] = usession
		if err := t.Execute(w, nil, dataMap); err != nil {
			log.Println(" - respnose-generator.go  View  : %s", err.Error())
		}
	}
}

// JSON Returns the data in form of "JSON" for the incoming
// request
func renderJSON(w http.ResponseWriter, data interface{}) {

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
